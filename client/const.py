import pyaudio

HOST = "10.0.0.218"
TCP_PORT = 9901
WS_PORT = 9902


FRAME_TYPE = 2  # ws frame type 1=text 2=binary
AUDIO_TYPE = 1
VIDEO_TYPE = 2

TCP_BUFSIZE = 4096


CHUNK = 512*2
FORMAT = pyaudio.paInt16
CHANNELS = 1
RATE = 48000
AUDIO_BUFSIZE = 1024


ROOM_ID = 10240
TOKEN1 = b'00000000000000000000000000000000'
TOKEN2 = b'00000000000000000000000000000001'

HEADER_LENGTH = 6


def get_user_header(i):
    header = bytes()
    loginDataType = 0
    h1 = FRAME_TYPE.to_bytes(length=1, byteorder='little', signed=False)
    h2 = loginDataType.to_bytes(length=1, byteorder='little', signed=False)
    rids = ROOM_ID.to_bytes(length=8, byteorder='little', signed=True)
    if i is None or i == 1:
        return (h1 + h2 + TOKEN1 + rids)

    return (h1 + h2 + TOKEN2 + rids)


def get_audio_header(size):
    h1 = FRAME_TYPE.to_bytes(length=1, byteorder='little', signed=False)
    h2 = AUDIO_TYPE.to_bytes(length=1, byteorder='little', signed=False)
    sizes = size.to_bytes(length=4, byteorder="little", signed=False)
    return (h1 + h2 + sizes)


def get_video_header(size):
    h1 = FRAME_TYPE.to_bytes(length=1, byteorder='little', signed=False)
    h2 = VIDEO_TYPE.to_bytes(length=1, byteorder='little', signed=False)
    sizes = size.to_bytes(length=4, byteorder="little", signed=False)
    return (h1 + h2 + sizes)
