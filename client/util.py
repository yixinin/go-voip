

def get_body_length(header):
    return int.from_bytes(header[2:], byteorder="little", signed=False)


def get_data_type(header):
    return header[1]
