import socket
import logging

from common.bingo import Bingo


class Server:
    def __init__(self, port, listen_backlog, endingMessage):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._endingMessage = endingMessage

    def run(self, statuses: dict):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not statuses['killWasCalled']:
            client_sock = self.__accept_new_connection()
            if statuses['killWasCalled']:
                break
            self.__handle_client_connection(client_sock)
        logging.info('server stopped listening for connections, will shut down in short time.')

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            msg = self.getMessage(client_sock)
            if msg == 'test':
                client_sock.send("{}\n".format(msg).encode('utf-8'))
                return
            addr = client_sock.getpeername()
            bingoService = Bingo(addr[1])
            bingoService.processMessage(msg)
            self.sendMessage(client_sock, msg)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def sendMessage(self, clientSock, msg):
        eightKb = 1024*8
        finalMessage = msg + self._endingMessage
        for i in range(0, len(finalMessage), eightKb):
            bytesSent = clientSock.send(finalMessage[i: i+eightKb].encode('utf-8'))
            i -= eightKb - bytesSent

    def getMessage(self, clientSock):
        eightKb = 1024*8
        msg = ''
        while True:
            msg += clientSock.recv(eightKb).rstrip().decode('utf-8')
            if msg == 'test' or msg.endswith(self._endingMessage):
                break
        return msg.rstrip(self._endingMessage)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def stopConnection(self):
        """
        stop listening for connections
        """
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.connect(self._server_socket.getsockname())
