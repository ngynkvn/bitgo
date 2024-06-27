from ulid import ULID
from socket import socket
from dataclasses import dataclass, asdict
import json


class Request:
    id: str = ""

    def __init__(self):
        super().__init__()
        self.id = str(ULID())


@dataclass
class RequestStatus(Request):
    method: str = "status"
    id: str = ""


@dataclass
class RequestVersion(Request):
    method: str = "version"
    id: str = ""


@dataclass
class RequestAddTorrent(Request):
    path: str
    method: str = "add"


@dataclass
class Response:
    result: str
    error: str
    id: str


class BackendClient:
    def __init__(self, socket: socket):
        self.socket = socket

    def send_status(self) -> Response:
        return self.send(RequestStatus())

    def send_version(self) -> Response:
        return self.send(RequestVersion())

    def send_add_torrent(self, path: str) -> Response:
        return self.send(RequestAddTorrent(path))

    def send(self, req) -> Response:
        ser = json.dumps(asdict(req))
        self.socket.sendall(bytes(ser, "utf-8"))
        data = self.socket.recv(2048).decode("utf-8")
        data = json.loads(data)
        return Response(**data)
