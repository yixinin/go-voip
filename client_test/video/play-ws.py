import websocket
import pyaudio
import wave
import numpy as np
import cv2

frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video


def Recv(ws):
    preBuf = bytes()
    while 1:
        buf = (preBuf + ws.recv())
        header = buf[:2+4]
        body = buf[2+4:]

        length = int.from_bytes(header[2:], byteorder="little", signed=False)
        read = body.__len__()  # 已读取的长度

        if read == length:
            play(body)
            if cv2.waitKey(1) & 0xFF == ord('q'):
                break
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
            if cv2.waitKey(1) & 0xFF == ord('q'):
                break

        elif read > length:
            # 粘包  分解
            while read > length:
                play(body[:length])
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break
                body = body[length:]
                read -= length
            preBuf = body

        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

    cv2.destroyAllWindows()


def play(data):
    arr = np.frombuffer(data, np.uint8)
    frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)

    cv2.imshow("video", frame)


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

    # ws.close()
