import yaml
import sys


def addClients(amount, services: dict):
    keys = services.keys()
    copyKeys = []
    for k in keys:
        if k.startswith('client'):
            copyKeys.append(k)
    for k in copyKeys:
        del services[k]
    for i in range(amount):
        clientId = i + 1
        services[f'client{clientId}'] = {
            'container_name': f'client{clientId}',
            'image': 'client:latest',
            'entrypoint': '/client',
            'environment':
                [f'CLI_ID={clientId}', 'CLI_LOG_LEVEL=DEBUG'],
            'networks': ['testing_net'],
            'depends_on': ['server']
        }


if __name__ == "__main__":
    amountToAdd = sys.argv[1]
    if not (amountToAdd and amountToAdd.isdigit()):
        print('la cagaste pibe')
        exit(1)
    amountToAdd = int(amountToAdd)
    info = {}
    with open('docker-compose-dev.yaml', 'r') as f:
        info = yaml.load(f, Loader=yaml.FullLoader)
    addClients(int(amountToAdd), info['services'])
    with open('docker-compose-dev.yaml', 'w') as f:
        yaml.dump(info, f)
