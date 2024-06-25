from socket import socket, AF_UNIX, SOCK_STREAM
from textual import events
from textual.screen import Screen
from textual.suggester import SuggestFromList
from textual.app import App, ComposeResult
from textual.widgets import Header, Footer, Log, Input, ListView

from protocol import BackendClient

SOCKET_FILE_PATH = "/tmp/bitgo.sock"


class AddTorrent(Screen):
    BINDINGS = [("escape", "app.pop_screen", "Pop screen")]

    def on_mount(self):
        self.title = "Add torrent"

    def on_key(self, event: events.Key):
        pass

    def compose(self) -> ComposeResult:
        yield Header(show_clock=True)
        yield Input(
            placeholder="path",
            suggester=SuggestFromList(suggestions=["./torrent/debian.torrent"]),
        )
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

    def on_key(self, event: events.Key) -> None:
        if event.key == "A":
            print("A")

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
