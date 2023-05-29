const net = require('net');

const FrameTypeCode = {
    close: 0,
    ping: 1,
    pong: 2,
    assert: 3,
    produce: 4,
    consume: 5,
    ack: 6,
    delegate: 7
};

const FrameCodeType = {
    0: 'close',
    1: 'ping',
    2: 'pong',
    3: 'assert',
    4: 'produce',
    5: 'consume',
    6: 'ack',
    7: 'delegate'
};

class Frame {
    /**
     * @param {keyof FrameTypeCode} type
     * @param {string | undefined} subject
     * @param {string | undefined} body
     */
    constructor(type, subject, body) {
        this.code = FrameTypeCode[type];
        this.subject = subject || '';
        this.body = body || '';
    }

    encode() {
        return `<${this.code}:${this.subject}:${this.body}>`;
    }

    /**
     * @param {Buffer} v
     */
    static decode(v) {
        const str = v.toString();
        const code = +str[1];
        let subject = '';
        let body = '';

        let i = 3;

        while (i < str.length) {
            if (str.at(i) === ':') break;
            subject += str.at(i);
            i++;
        }

        i++;

        while (i < v.length) {
            if (str.at(i) === '>') break;
            body += str.at(i);
            i++;
        }

        return new Frame(FrameCodeType[code], subject, body);
    }
}

class Client {
    constructor() {
        this._socket = new net.Socket();
        this._socket.connect(6789, '127.0.0.1', () => {
            this._interval = setInterval(this.ping.bind(this), 10000);
        });

        this._socket.on('data', buf => {
            console.info(Frame.decode(buf));
        });

        this._socket.on('end', () => {
            this._socket.destroy();
            clearInterval(this._interval);
        });
    }

    ping() {
        this._socket.write(new Frame('ping').encode());
    }

    close() {
        this._socket.write(new Frame('close').encode());
    }

    /**
     * @param {string | RegExp} key
     */
    assert(key) {
        this._socket.write(new Frame('assert', key).encode());
    }

    /**
     * @param {string | RegExp} key
     * @param {string} payload
     */
    produce(key, payload) {
        this._socket.write(new Frame('produce', key, payload).encode());
    }

    /**
     * @param {string | RegExp} key
     */
    consume(key) {
        this._socket.write(new Frame('consume', key).encode());
    }
}

new Client();
