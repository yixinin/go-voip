from socket import *
import pyaudio
import wave
import numpy as np

frameType = 2  # ws frame type 1=text 2=binary
dataType = 1   # live data type 1=audio 2=video
BUFSIZE = 4096


def Recv(tcpClient):
    preBuf = bytes()
    while 1:
        buf = (preBuf + tcpClient.recv(BUFSIZE))
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
                subBody = tcpClient.recv(BUFSIZE)

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


def conn():
    # 发送登录信息
    header = bytes()
    h1 = frameType.to_bytes(length=1, byteorder='little', signed=True)
    h2 = dataType.to_bytes(length=1, byteorder='little', signed=True)

    token = b'00000000000000000000000000000001'
    rid = 10240
    x2 = rid.to_bytes(length=8, byteorder='little', signed=True)
    header = bytes(h1+h2+token+x2)

    HOST = '10.0.0.23'
    PORT = 9901
    ADDR = (HOST, PORT)

    tcpClient = socket(AF_INET, SOCK_STREAM)
    tcpClient.connect(ADDR)
    tcpClient.send(header)
    return tcpClient


def play(data):
    arr = np.frombuffer(data, np.uint8)
    frame = cv2.imdecode(arr, cv2.IMREAD_COLOR)

    cv2.imshow("video", frame)


if __name__ == "__main__":
    tcpClient = conn()
    Recv(tcpClient)

    tcpClient.close()
