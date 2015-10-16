import os
import shutil

from charmhelpers.core import hookenv, host
from charmhelpers.core.templating import render
from charms.reactive import hook, when, when_not, is_state, set_state, remove_state
from charms.reactive.bus import get_states


@hook('install')
def install():
    if is_state('oostore.available'):
        return
    host.adduser('oostore', system_user=True)
    host.add_group('oostore', system_group=True)
    host.add_user_to_group('oostore', 'oostore')
    install_workload()
    set_state('oostore.available')


@hook('upgrade-charm')
def upgrade():
    if is_state('oostore.started'):
        host.service_stop('oostore')
    install_workload()
    if is_state('oostore.started'):
        host.service_start('oostore')


def install_workload():
    config = hookenv.config()
    deployment = config['deployment']
    if not os.path.exists('/srv/oostore/%s/bin' % (deployment)):
        os.makedirs('/srv/oostore/%s/bin' % (deployment), mode=0o755)
    shutil.copyfile('files/oostore', '/srv/oostore/%s/bin/oostore' % (deployment))
    os.chmod('/srv/oostore/%s/bin/oostore' % (deployment), 0o755)


@hook('config-changed')
def config_changed():
    config = hookenv.config()
    if config.changed('http_port'):
        if config.previous('http_port'):
            hookenv.close_port(config.previous('http_port'))
        hookenv.open_port(config['http_port'])
    if config.changed('https_port'):
        if config.previous('https_port'):
            hookenv.close_port(config.previous('https_port'))
        hookenv.open_port(config['https_port'])
    set_state('oostore.configured')
 

@when('oostore.start')
@when_not('oostore.started')
def start_oostore():
    host.service_start('oostore')
    set_state('oostore.started')


@when('oostore.started')
@when_not('oostore.start')
def stop_oostore():
    host.service_stop('oostore')
    remove_state('oostore.started')


@when('oostore.configured', 'database.connected', 'database.database.available')
def setup(pg, _):
    config = hookenv.config()
    render(source="upstart",
        target="/etc/init/oostore.conf",
        owner="root",
        perms=0o644,
        context={
            'cfg': config,
            'db': pg,
        })

    deployment = config['deployment']
    if not os.path.exists('/srv/oostore/%s/etc' % (deployment)):
        os.makedirs('/srv/oostore/%s/etc' % (deployment), mode=0o755)

    cert_file = '/srv/oostore/%s/etc/cert.pem' % (deployment)
    if config.get('cert'):
        with open(cert_file, "w") as fh:
            fh.write(config['cert'])
    else:
        if os.path.exists(cert_file):
            os.unlink(cert_file)

    key_file = '/srv/oostore/%s/etc/key.pem' % (deployment)
    if config.get('key'):
        with open(key_file, "w") as fh:
            fh.write(config['key'])
        os.chmod(key_file, 0o600)
    else:
        if os.path.exists(cert_file):
            os.unlink(key_file)

    set_state('oostore.start')
    hookenv.status_set('maintenance', 'Starting oostore')


@when('oostore.available')
@when_not('database.connected')
def missing_db():
    hookenv.log("%s" % (str(get_states())))
    remove_state('oostore.start')
    hookenv.status_set('blocked', 'Please add relation to postgresql')


@when('database.connected')
@when_not('database.database.available')
def waiting_db(pg):
    hookenv.log("%s" % (str(get_states())))
    remove_state('oostore.start')
    hookenv.status_set('waiting', 'Waiting for postgresql')


@when('oostore.started')
def oostore_started():
    config = hookenv.config()
    hookenv.open_port(config['http_port'])
    if config.get('cert') and config.get('key'):
        hookenv.open_port(config['https_port'])
    
    hookenv.status_set('active', 'Ready')


@when('website.available')
def configure_website(website):
    config = hookenv.config()
    website.configure(config['port'])
