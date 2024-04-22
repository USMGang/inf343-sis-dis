# Integrantes

Sebastian Gutierrez Milla - 202073601-4

# Instrucciones

El laboratorio esta implementado con makefile, estos serian los comandos:
- `make server`: Sirve para ejecutar la parte del server alojada en el contenedor `lab2_super_tierra` de docker.
- `make client`: Sirve para ejecutar la parte del cliente.

Estos 2 comandos ejecutan pasos intermedios, los cuales son:
- `make proto`: Crea los archivos necesarios para implementar gRPC
- `make docker`: Crea el contenedor de docker

# Comentarios

Por asuntos democráticos, la función `solicitarM(ID, AT, MP)` sufrió un cambio de
nombre a `Resupply(ID, AT, MP)`.

Todos los comandos de docker incluyen `sudo` en caso de que se necesite

Se puede ejecutar `make serverd` para ejecutar el contenedor de docker en modo
detached.

También se implemento `make s` en caso de que solamente se requiera hacer un
`docker run` y no volver a construir el contenedor.
