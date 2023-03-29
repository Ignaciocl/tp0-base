#!/usr/bin/env python3
import signal
from configparser import ConfigParser
from common.server import Server
import logging
import os

statuses = {'killWasCalled': False}


def initialize_config():
    """ Parse env variables or config file to find program config params

    Function that search and parse program configuration parameters in the
    program environment variables first and the in a config file.
    If at least one of the config parameters is not found a KeyError exception
    is thrown. If a parameter could not be parsed, a ValueError is thrown.
    If parsing succeeded, the function returns a ConfigParser object
    with config parameters
    """

    config = ConfigParser(os.environ)
    # If config.ini does not exists original config object is not modified
    config.read("config.ini")

    config_params = {}
    try:
        config_params["port"] = int(os.getenv('SERVER_PORT', config["DEFAULT"]["SERVER_PORT"]))
        config_params["listen_backlog"] = int(
            os.getenv('SERVER_LISTEN_BACKLOG', config["DEFAULT"]["SERVER_LISTEN_BACKLOG"]))
        config_params["logging_level"] = os.getenv('LOGGING_LEVEL', config["DEFAULT"]["LOGGING_LEVEL"])
        config_params["endingMessage"] = os.getenv('ENDING_MESSAGE', config["DEFAULT"]["ENDING_MESSAGE"])
        config_params["batchMessageEnd"] = os.getenv('ENDING_BATCH', config["DEFAULT"]["ENDING_BATCH"])
        config_params["maxConnections"] = os.getenv('MAX_CONNECTIONS_ALIVE', config["DEFAULT"]["MAX_CONNECTIONS_ALIVE"])
    except KeyError as e:
        raise KeyError("Key was not found. Error: {} .Aborting server".format(e))
    except ValueError as e:
        raise ValueError("Key could not be parsed. Error: {}. Aborting server".format(e))

    return config_params


def main():
    config_params = initialize_config()
    logging_level = config_params["logging_level"]
    port = config_params["port"]
    listen_backlog = config_params["listen_backlog"]
    endingMessage = config_params["endingMessage"]
    batchMessageEnd = config_params["batchMessageEnd"]
    maxConnections = int(config_params["maxConnections"])

    initialize_log(logging_level)

    # Log config parameters at the beginning of the program to verify the configuration
    # of the component
    logging.debug(f"action: config | result: success | port: {port} | "
                  f"listen_backlog: {listen_backlog} | logging_level: {logging_level}")

    # Initialize server and start server loop
    server = Server(port, listen_backlog, endingMessage, batchMessageEnd)

    def killWasCalled(*args):
        logging.info('sigterm was called. Shutting down')
        statuses['killWasCalled'] = True
        server.stopConnection()

    signal.signal(signal.SIGTERM, killWasCalled)
    server.run(statuses, maxConnections)


def initialize_log(logging_level):
    """
    Python custom logging initialization

    Current timestamp is added to be able to identify in docker
    compose logs the date when the log has arrived
    """
    logging.basicConfig(
        format='%(asctime)s %(levelname)-8s %(message)s',
        level=logging_level,
        datefmt='%Y-%m-%d %H:%M:%S',
    )


if __name__ == "__main__":
    main()
