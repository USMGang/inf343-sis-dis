# Datos
Sebastian Gutierrez 202073601-4
Juan Diego Arias - 202073590-5

# Instrucciones

El laboratorio esta implementado con makefile:
- `make player`: Ejecuta el código de mercenario controlado por el usuario.
- `make bot`: Ejecuta el código de mercenario controlado por un bot.

- `make proto`: Compila todos los archivos con proto.

- `make director`: Ejecuta el contenedor del director
- `make doshbank`: Ejecuta el contenedor del doshbank
- `make namenode`: Ejecuta el contenedor del namenode
- `make datanode ID=`: Ejecuta el contenedor del datanode según el id que se le
entregue

- `make clear`: Elimina todos los txt y limpia otros archivos (binarios si se
ejecutaba en local)

También implementamos un .env para colocar los hosts. 

Se asume que las maquinas virtuales ejecutaran:
    - dist069: make director, make doshbank, make namenode, make player
    - dist070: make datanode ID=1 make bot make bot make bot
    - dist071: make datanode ID=2 make bot make bot
    - dist071: make datanode ID=3 make bot make bot

Según los host en .env

# Notas

- Se asume que todos los datos ingresados en las interfaces serán validos.
- Se editó un poco la lógica del piso 3 y en vez de que el numero del patriarca y del mercenario fuera igual ahora es menor-igual.
- El director no termina cuando se acaba la misión.
- No se incluyó la lógica de poder hacer que el director pueda ver las opciones
  desde la interfaz, pero la estructura esta hecha, ademas de los .txt
