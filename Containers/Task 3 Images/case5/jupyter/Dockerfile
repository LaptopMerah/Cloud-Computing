FROM alpine:3.18

RUN apk update && apk add --no-cache \
    python3 \
    py3-pip \
    supervisor \
    curl \
    bash \
    gcc \
    python3-dev \
    musl-dev \
    linux-headers \
    && rm -rf /var/cache/apk/*

RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir -r /tmp/requirements.txt

RUN adduser -D -h /home/jupyter jupyter
RUN mkdir -p /home/jupyter/notebooks && chown -R jupyter:jupyter /home/jupyter

COPY start.sh /tmp/start.sh
COPY supervisord.conf /tmp/supervisord.conf
COPY jupyter_config.py /home/jupyter/.jupyter/jupyter_notebook_config.py

RUN chmod +x /tmp/start.sh
RUN chown -R jupyter:jupyter /home/jupyter

EXPOSE 8888

WORKDIR /home/jupyter/notebooks

VOLUME ["/home/jupyter/notebooks"]

ENTRYPOINT ["/bin/sh", "/tmp/start.sh"]
