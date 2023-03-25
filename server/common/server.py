import json
import socket
import logging

from common.bingo import Bingo, LotteryManager


class Server:
    def __init__(self, port, listen_backlog, endingMessage, endingBatch):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._endingMessage = endingMessage
        self._endingBatch = endingBatch
        self._statusWasKilled = False

    def run(self, statuses: dict):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not statuses['killWasCalled']:
            client_sock = self.__accept_new_connection()
            self.__handle_client_connection(client_sock)
        logging.info('server stopped listening for connections, will shut down in short time.')
        self._statusWasKilled = True

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        processed = 0
        try:
            while not self._statusWasKilled:
                msg, keepProcessing = self.getMessage(client_sock)
                if msg == 'test':
                    client_sock.send("{}\n".format(msg).encode('utf-8'))
                    return
                addr = client_sock.getpeername()
                bingoService = Bingo(addr[1])
                processedThisIter: dict = bingoService.processMessage(msg)
                self.sendMessage(client_sock, json.dumps({"amount_processed": processedThisIter.get('amountProcessed', 0), "status": processedThisIter.get('status', "allOgre"), "winners": processedThisIter.get('winners')}))
                processed += processedThisIter.get('amountProcessed', 0)
                if not keepProcessing:
                    lm = LotteryManager()
                    lm.agencyFinished(addr[1])
                    break
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e} | amount processed = {processed}")
        finally:
            logging.info(f"action: finish_processing | result: ok | amountProcessed = {processed}")
            client_sock.close()

    def sendMessage(self, clientSock, msg):
        eightKb = 1024*8
        finalMessage = msg + self._endingMessage
        for i in range(0, len(finalMessage), eightKb):
            bytesSent = clientSock.send(finalMessage[i: i+eightKb].encode('utf-8'))
            i -= eightKb - bytesSent

    def getMessage(self, clientSock):
        """
        :param clientSock: Sock of the client to receive the message
        :returns: msg to process and if the loop should continue or not depending on how the message was received
        """
        eightKb = 1024*8
        msg = ''
        while True:
            msg += clientSock.recv(eightKb).rstrip().decode('utf-8')
            if msg == 'test':
                return msg, False
            if msg.endswith(self._endingMessage) or msg.endswith(self._endingBatch):
                break
        return msg.rstrip(self._endingMessage).rstrip(self._endingBatch), not msg.endswith(self._endingBatch)

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
