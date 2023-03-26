# TP0: Docker + Comunicaciones + Concurrencia
El tp esta separado por los ejercicios pedidos a resolver, para ver cada uno se tiene que acceder a cada uno de los branchs en especifico. Los branches se relacionan con los ejercicios de la siguiente manera:
 - Ejercicio 1: [ex1](https://github.com/Ignaciocl/tp0-base/tree/ex1)
 - Ejercicio 2: [ex2](https://github.com/Ignaciocl/tp0-base/tree/ex2)
 - Ejercicio 3: [ex3](https://github.com/Ignaciocl/tp0-base/tree/ex3)
 - Ejercicio 4: [ex4](https://github.com/Ignaciocl/tp0-base/tree/ex4)
 - Ejercicio 5: [ej5](https://github.com/Ignaciocl/tp0-base/tree/ej5)
 - Ejercicio 6: [ej6](https://github.com/Ignaciocl/tp0-base/tree/ej6)
 - Ejercicio 7: [ej7](https://github.com/Ignaciocl/tp0-base/tree/ej7)
 - Ejercicio 8: [ej8](https://github.com/Ignaciocl/tp0-base/tree/ej8)


## Instrucciones de uso

## Parte 1: Introducción a Docker
En esta primera parte se explicara como ejecutar cada punto que no implique el uso de make docker-compose-up, en caso de que no se especifique un ejercicio se presupone que se utilizara make docker-compose-up y que eso es suficiente
### Ejercicio N°1:
Para la ejecucion de este ejercicio lo unico que se debe hacer es:
```
python3 addNetworks.py N
```
Siendo N un numero valido de clientes, en caso de que el numero sea negativo no se agregara ningun cliente, en caso de que no sea numero se levantara un error y el codigo no se ejecutara.

### Ejercicio N°3:
Este se ejecutara de la misma forma que se dijo en un comienzo, para detectar el correcto funcionamiento se puede acceder a la maquina levantada con docker y ahi se vera el correcto funcionamiento del mismo.

## Parte 2: Repaso de Comunicaciones

En esta seccion se utilizo un protocolo en el cual se mandan los mensajes en formato json, seguidos de una cadena de caracteres configurables por entorno para determinar cuando el mensaje termino de ser enviado.
Se evita los fenomenos de short write y short read contando los caracteres enviados y enviando los que hacen falta devuelta en la siguiente iteracion. Se envian por paquete 8kb, lo cual para la informacion enviada tiene capacidad de sobra.
Para evitar el short read lo que se hace es esperar hasta el mensaje final, entonces ahi no se tienen en cuanto cuantos caracteres se leyeron sino saber que se termino de leer correctamente.
Para mantener backwards compatibility se hizo una exepcion en el servidor para cuando recibe el mensaje "test" provisto por el script de nc, asi nc puede seguir funcionando correctamente hasta en estos puntos.

### Ejercicio N°6:
En este ejercicio se hizo una modificacion y se distinguio entre dos tipos de envios de mensaje:
 - Una parte de un batch, en ese caso terminara con un mensaje que le indica al servidor que seguira enviando.
 - El batch final, le indica al servidor que termino de enviar entonces puede cerrar la conexion.
El ejercicio en si era ambiguo sobre si se queria que se enviaran batchs en nuevas conexiones, por lo que me decidi en lo que ahorraria tiempo de coneccion a costa de bloquear la coneccion para otros clientes que quieran subir sus datos. Cambiar esto para el otro caso es trivial, pero es un punto importante a notar.

### Ejercicio N°7:
Dado el approach que se tomo en el ejercicio numero 6, la primera parte fue trivial ya que el servidor ya se enteraba cuando terminaba el envio de batchs. En la segunda parte fue implementada en este ejercicio.
## Parte 3: Repaso de Concurrencia
Se utilizo la libreria de multiprocessing de python, la cual permite la creacion de nuevos procesos, el hecho de uso de pools parecio muy excesivo para el alcance de este trabajo, por lo que a cada nueva coneccion el servido inicia un nuevo proceso que despues cuando se mande la señal de terminar, joineara todos los procesos.
Para sincronizar el acceso a las secciones criticas se utilizo el recurso de locks y Values (este ultimo siendo un value que se comparte en memoria y que tiene un lock interno para manejar secciones criticas a pedido).
Se utilizan los locks para impedir que dos procesos escriban al mismo tiempo, y se utilizan una mezcla de los values y locks para impedir que un proceso lea mientras otros escriben, pero que procesos que lean no se impidan el paso unos a otros.
