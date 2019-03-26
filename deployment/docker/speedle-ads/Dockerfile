FROM oraclelinux:7-slim
WORKDIR /speedle
COPY ./* /speedle/
EXPOSE 6734
ENV PATH="/speedle:${PATH}"

RUN mkdir -p /var/lib/speedle; \
    echo "{}" > /var/lib/speedle/policies.json

ENTRYPOINT ["entrypoint.sh"]

