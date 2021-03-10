FROM brunoanjos/nova-server-base:latest

ENV executable="executable"
COPY $executable .
COPY store_items.json .

CMD ["sh", "-c", "./$executable"]