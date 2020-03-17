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

    preBuf = bytes()
    while 1:
        buf = (preBuf + ws.recv())
        header = buf[:2+4]
        body = buf[2+4:]

        length = int.from_bytes(header[2:], byteorder="little", signed=False)
        read = body.__len__()  # 已读取的长度

        if read == length:
            play(body)

        elif read < length:
            # 拆包 合并
            while read < length:
                unRead = length - read
                subBody = ws.recv()

                read += subBody.__len__()
                if read >= length:
                    body = (body+subBody[:length-read])
                    preBuf = subBody[length-read:]
                    read == length
                else:
                    body = (body+subBody)
            play(body)

        elif read > length:
            # 粘包  分解
            while read > length:
                play(body[:length])

                body = body[length:]
                read -= length
            preBuf = body

    # 停止数据流
    stream.stop_stream()
    stream.close()

    # 关闭 PyAudio
    p.terminate()


def play(stream, body):
    stream.write(body)


def conn():
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
    return ws


if __name__ == "__main__":
    ws = conn()
    Recv(ws)
