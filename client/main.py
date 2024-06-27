from socket import socket, AF_UNIX, SOCK_STREAM
from typing import Any
from textual import events
from textual.screen import Screen
from textual.suggester import SuggestFromList
from textual.message import Message
from textual.app import App, ComposeResult
from textual.widgets import Header, Footer, Log, Input, ListView, Button

from protocol import BackendClient

SOCKET_FILE_PATH = "/tmp/bitgo.sock"


class AddTorrent(Screen):
    BINDINGS = [("escape", "app.pop_screen", "Pop screen")]

    def on_mount(self):
        self.title = "Add torrent"

    def on_key(self, event: events.Key):
        pass

    class Form(Message):
        def __init__(self, path: str):
            super().__init__()
            self.path = path

    def on_button_pressed(self, event: Button.Pressed) -> None:
        self.post_message(self.Form(self.query_one(Input).value))
        self.app.pop_screen()

    def compose(self) -> ComposeResult:
        yield Header(show_clock=True)
        yield Input(
            placeholder="path",
            suggester=SuggestFromList(suggestions=["./torrent/debian.torrent"]),
        )
        yield Button("Add", variant="primary")
        yield Log()
        yield Footer()


class Main(App):
    SCREENS = {"addtorrent": AddTorrent()}
    BINDINGS = [("a", "push_screen('addtorrent')", "Add a torrent")]

    def __init__(self, s: socket):
        super().__init__()
        self.backend = BackendClient(s)
        self.SOCKET = s

    def on_mount(self):
        resp = self.backend.send_version()

        self.version = resp.result
        self.title = self.version

    def on_add_torrent_form(self, message: AddTorrent.Form) -> None:
        resp = self.backend.send_add_torrent(message.path)
        self.query_one(Log).write_line(resp.result)

    def on_ready(self):
        pass

    def compose(self):
        yield Header(show_clock=True)
        yield ListView()
        yield Log()
        yield Footer()


if __name__ == "__main__":
    with socket(AF_UNIX, SOCK_STREAM) as sock:
        sock.connect(SOCKET_FILE_PATH)
        app = Main(sock)
        app.run()
