# from websocket import create_connection
import websocket
import pyaudio
import wave
import numpy as np

frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video
BUFSIZE = 1024+2+8


def Recv(ws):
    CHUNK = 512*2
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 48000

    p = pyaudio.PyAudio()
    stream = p.open(rate=RATE, channels=CHANNELS, format=FORMAT, output=True)

    data = bytes()
    while 1:
        buf = (data + ws.recv())
        BUFSIZE = 12 + \
            int.from_bytes(buf[2:4], byteorder="little", signed=False)
        siz = buf.__len__()
        while siz >= BUFSIZE:
            stream.write(buf[8+4:BUFSIZE])
            buf = buf[BUFSIZE:]
            siz = buf.__len__()
        data = buf

    # 停止数据流
    stream.stop_stream()
    stream.close()

    # 关闭 PyAudio
    p.terminate()


if __name__ == "__main__":
    # 发送登录信息
    header = bytes()
    h1 = frameType.to_bytes(length=1, byteorder='little', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='little', signed=True)

    token = b'00000000000000000000000000000001'
    rid = 10240
    x2 = rid.to_bytes(length=8, byteorder='little', signed=True)
    header = bytes(h1+h2+token+x2)

    addr = "ws://127.0.0.1:9902/live"

    # websocket.enableTrace(True)
    ws = websocket.create_connection(addr)
    ws.send_binary(header)

    Recv(ws)
