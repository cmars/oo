import shutil

from charmhelpers.core import hookenv, host
from charmhelpers.core.templating import render
from charms.reactive import hook, when, when_not, is_state, set_state, remove_state


@hook('install')
def install():
    if is_state('oostore.available'):
        return
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
    shutil.copyfile('files/oostore', '/srv/oostore/%s/bin/oostore' % (deployment))


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
    set_state('oostore.configured', config)
 

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


@when('oostore.configured', 'database.available')
def setup(config, pg):
    render(source="templates/upstart",
        target="/etc/init/oostore.conf",
        owner="root",
        perms=0o644,
        context={
            cfg: config,
            db: pg,
        })

    cert_file = '/srv/oostore/%s/etc/cert.pem' % (config['deployment'])
    if config.get('cert'):
        with open(cert_file, "w") as fh:
            fh.write(config['cert'])
    else:
        os.unlink(cert_file)

    key_file = '/srv/oostore/%s/etc/key.pem' % (config['deployment'])
    if config.get('key'):
        with open(key_file, "w") as fh:
            fh.write(config['key'])
        os.chmod(key_file, 0o600)
    else:
        os.unlink(key_file)

    set_state('oostore.start')
    status_set('maintenance', 'Starting oostore')


@when_not('database.connected')
def missing_db():
    remove_state('oostore.start')
    hookenv.status_set('blocked', 'Please add relation to postgresql')


@when('database.connected')
@when_not('database.available')
def waiting_db():
    remove_state('oostore.start')
    hookenv.status_set('waiting', 'Waiting for postgresql')


@when('oostore.started')
def oostore_started():
    hookenv.status_set('active', 'Ready')
