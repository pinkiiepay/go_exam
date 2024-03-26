# Exámen Golang

Este es un proyecto para un exámen de Go tomando las siguientes consideraciones:
Realiza el siguiente ejercicio, comparte la liga del repositorio de tu prueba. 

1. Crear un servicio de registro de usuario que reciba como parámetros usuario, correo,
teléfono y contraseña.

2. El servicio deberá validar que el correo y telefono no se encuentren registrados, de lo
contrario deberá retornar un mensaje “el correo/telefono ya se encuentra registrado”.

3. Deberá validar que la contraseña sea de 6 caracteres mínimo y 12 máximo y contener
al menos una mayúscula, una minúscula, un carácter especial (@ $ o &) y un número.

4. Validar que el teléfono sea a 10 dígitos y el correo tenga un formato válido.

5. Crear un servicio login que reciba como parámetros usuario o correo y contraseña.

6. El servicio debe devolver un token jwt.

7. Deberá validar que el usuario o correo y contraseña sean válidos, de lo contrario
retorna un mensaje “usuario / contraseña incorrectos”.

8. En ambos servicios se deberá validar que todos los parámetros solicitados vayan en el
cuerpo de la petición, de lo contrario retorna un mensaje con el campo faltante.




## Notas

Para que funciones se deben descargar las siguientes librerías:

	"github.com/golang-jwt/jwt"
	"github.com/microsoft/go-mssqldb"
	"gopkg.in/go-playground/validator.v9"

con el comando go get.

Se utilizó una base de datos en SQL SERVER local montada en docker:

https://learn.microsoft.com/es-es/sql/linux/quickstart-install-connect-docker?view=sql-server-ver16&tabs=cli&pivots=cs1-cmd


Y se implementó las siguientes líneas de código:

```
CREATE schema EXAM_GO

CREATE TABLE EXAM_GO.Users (id INT NOT NULL IDENTITY, user_name NVARCHAR(50), email NVARCHAR(50), phone NVARCHAR(10), password NVARCHAR(12), PRIMARY KEY (id));
```