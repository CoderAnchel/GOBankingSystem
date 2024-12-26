# GO Banking System
![banking ada](https://github.com/user-attachments/assets/f7f8853f-2a74-4794-b3ff-fed06399c34c)

> Sistema bancario desarrollada en Go (Golang) que proporciona funcionalidades bancarias y de trading, utilizando una arquitectura basada en REST API para manejar las interacciones entre los usuarios y el sistema.


## 🎯 Tareas
1. **Tarea 1**: Acciones de Usuario
2. **Tarea 2**: Restablecimiento de Contraseña y OTP
3. **Tarea 3**: Creación y Gestión de PIN
4. **Tarea 4**: Transacciones de Cuenta
5. **Tarea 5**: Operaciones de Mercado
6. **Tarea 6**: Seguridad
7. **Tarea 7**: Manejo de Errores
8. **Tarea 8**: Suscripciones y Bot de Trading

### 📑 Información detallada sobre las tareas

#### Tarea 1: Acciones de Usuario

Esta tarea se centra en acciones básicas relacionadas con el usuario, como registrar un nuevo usuario, iniciar sesión, recuperar detalles del usuario y de la cuenta, y cerrar sesión. Para estas acciones, necesitarás interactuar con varios endpoints, algunos de los cuales requieren autenticación.

- **Registro de Usuario**: Implementa la funcionalidad para registrar un usuario enviando la información requerida como nombre, correo electrónico, número de teléfono y contraseña. Este registro debe devolver el número de cuenta, que se utilizará para futuras operaciones.
    Cuerpo de la solicitud:
    ```json
    {
        "name":"Nuwe Test",
        "password":"NuweTest1$",
        "email":"nuwe@nuwe.com",
        "address":"Main St",
        "phoneNumber":"666888116"
    }
    ```
    Respuesta:
    ```json
    {
        "name": "Nuwe Test",
        "email": "nuwe@nuwe.com",
        "phoneNumber": "666888116",
        "address": "Main St",
        "accountNumber": "19b332",
        "hashedPassword": "$2a$10$vYWBxACqEIPeoT0O5b0faOHp4ITAHSBvoHDzBePW7tPqzpvqKLi6G"
    }
    ```
    El número de cuenta debe ser creado y asignado a la cuenta automáticamente por la aplicación y ser un UUID.

    Las verificaciones deben incluir:
    - No campos vacíos.
    - El formato del correo electrónico debe ser válido.
    - Reglas de contraseña que se detallarán más adelante.
    - Verificar si el correo electrónico o el número de teléfono ya existen.


- **Inicio de Sesión de Usuario**: Implementa un mecanismo de inicio de sesión utilizando un correo electrónico o número de cuenta junto con una contraseña. Después de una autenticación exitosa, el sistema debe devolver un token JWT, que se utilizará para todos los endpoints protegidos.
    Cuerpo de la solicitud:
    ```json
    {
        "identifier":"nuwe@nuwe.com",
        "password":"NuweTest1$"
    }
    ```
    Respuesta:
    ```json
    {
        "token": "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxOWIzMzIiLCJpYXQiOjE3Mjk1NzEzNzUsImV4cCI6MTcyOTY1Nzc3NX0.6qLQi50B1StobsUusfxCSqLdKeKOYdBZ3qj5Lw5G9eAdqoV1Juz3jyh2xwWByG7iJtusrhYPb_I62ycptcH4MA"
    }
    ```

    Si el identificador es inválido, devuelve lo siguiente con el código de estado 400:

    ```
    Usuario no encontrado para el identificador dado: nuwee@nuwe.com
    ```

    Si la contraseña es inválida, devuelve lo siguiente con el código de estado 401:

    ```
    Credenciales incorrectas
    ```

- **Obtener Información del Usuario**: Una vez iniciado sesión, utiliza el token JWT para recuperar información detallada del usuario (por ejemplo, nombre, correo electrónico, número de cuenta).
    Respuesta:
    ```json
    {
        "name": "Nuwe Test",
        "email": "nuwee@nuwe.com",
        "phoneNumber": "666888116",
        "address": "Main St",
        "accountNumber": "19b332",
        "hashedPassword": "$2a$10$vYWBxACqEIPeoT0O5b0faOHp4ITAHSBvoHDzBePW7tPqzpvqKLi6G"
    }
    ```
- **Obtener Información de la Cuenta**: Recupera información de la cuenta, como el saldo de la cuenta. Debes estar conectado.
    Respuesta:
    ```json
    {
        "accountNumber": "19b332",
        "balance": 0.0
    }
    ```
- **Cerrar Sesión**: Implementa un sistema de cierre de sesión que invalide el token JWT, asegurando que los usuarios no puedan acceder a los endpoints protegidos.


#### Tarea 2: Restablecimiento de Contraseña y OTP

Esta tarea implica implementar la funcionalidad de restablecimiento de contraseña utilizando Contraseñas de Un Solo Uso (OTPs). Se prueba la capacidad de enviar OTPs por correo electrónico, verificarlos y restablecer la contraseña del usuario.

- **Enviar OTP**: Crea un mecanismo que envíe un OTP al correo electrónico registrado del usuario para fines de restablecimiento de contraseña.
    Cuerpo de la solicitud:
    ```json
    {
        "identifier":"nuwe@nuwe.com"
    }
    ```
    Respuesta:
    ```json
    {
        "message": "OTP enviado exitosamente a: nuwee@nuwe.com"
    }
    ```
    Debes enviar un correo electrónico al usuario con el mensaje:
    OTP:XXXXXX donde X son números

- **Verificar OTP**: Implementa la funcionalidad para verificar el OTP proporcionado por el usuario. Tras una verificación exitosa, se debe generar un token de restablecimiento.
    Cuerpo de la solicitud:
    ```json
    {
        "identifier":"nuwe@nuwe.com",
        "otp":"893392"
    }
    ```
    Respuesta:
    ```json
    {
        "passwordResetToken": "9b1ae6c5-c247-434f-a8e7-893b026db107"
    }
    ```

- **Restablecer Contraseña**: Después de verificar el OTP, el usuario puede restablecer su contraseña utilizando el token de restablecimiento.
    Cuerpo de la solicitud:
    ```json
    {
        "identifier":"nuwe@nuwe.com",
        "resetToken": "9b1ae6c5-c247-434f-a8e7-893b026db107",
        "newPassword": "PassTest1$"
    }
    ```
    Respuesta:
    ```json
    {
        "message": "Contraseña restablecida exitosamente"
    }
    ```

- **Probar Nueva Contraseña**: Asegúrate de que la nueva contraseña funcione iniciando sesión con las credenciales actualizadas.

#### Tarea 3: Creación y Gestión de PIN

Esta tarea se centra en crear, actualizar y verificar PINs para transacciones sensibles. Este PIN debe ser utilizado para todas las transacciones. Estos endpoints deben requerir autenticación JWT.

- **Crear PIN**: Implementa una funcionalidad para crear un PIN asociado con la cuenta del usuario. Este PIN se utilizará en transacciones como depósitos, retiros y transferencias.
    Cuerpo de la solicitud:
    ```json
    {
        "pin":"1810",
        "password":"PassTest1$"
    }
    ```
    Respuesta:
    ```json
    {
        "msg": "PIN creado exitosamente"
    }
    ```

- **Actualizar PIN**: Los usuarios deben tener la capacidad de actualizar su PIN existente proporcionando su PIN antiguo y la contraseña de la cuenta.
    Cuerpo de la solicitud:
    ```json
    {
        "oldPin":"1810",
        "password":"PassTest1$",
        "newPin": "1811"
    }
    ```
    Respuesta:
    ```json
    {
        "msg": "PIN actualizado exitosamente"
    }
    ```

- **Crear PIN para Otras Cuentas**: Prueba la creación de un PIN para otra cuenta para asegurar la funcionalidad en múltiples cuentas.

#### Tarea 4: Transacciones de Cuenta

Esta tarea implica implementar transacciones financieras básicas como depósitos, retiros y transferencias de fondos. Además, incluye la visualización del historial de transacciones.

Para cualquier transacción, se debe verificar que haya fondos suficientes. Si no hay fondos suficientes, se debe mostrar el mensaje: Saldo insuficiente

Todas las transacciones deben requerir un JWT.

- **Depositar Dinero**: Crea una funcionalidad que permita a los usuarios depositar dinero en su cuenta utilizando el PIN correcto.
    Cuerpo de la solicitud:
    ```json
    {
        "pin":"1811",
        "amount":"100000.0"
    }
    ```
    Respuesta:
    ```json
    {
        "msg": "Dinero depositado exitosamente"
    }
    ```
    Verifica que los fondos se hayan agregado correctamente a la base de datos.

- **Retirar Dinero**: Implementa un sistema donde los usuarios puedan retirar dinero de su cuenta utilizando su PIN.
    Cuerpo de la solicitud:
    ```json
    {
        "amount":20000.0,
        "pin":"1811"
    }
    ```
    Respuesta:
    ```json
    {
        "msg": "Dinero retirado exitosamente"
    }
    ```
    Verifica que los fondos se hayan retirado correctamente de la base de datos.

- **Transferencia de Fondos**: Habilita la transferencia de fondos de una cuenta a otra utilizando números de cuenta y verificación de PIN.
    Cuerpo de la solicitud:
    ```json
    {
        "amount": 1000.0,
        "pin": "1811",
        "targetAccountNumber": "ed9050"
    }
    ```
    Respuesta:
    ```json
    {
        "msg": "Fondos transferidos exitosamente"
    }
    ```

- **Historial de Transacciones**: Implementa una función que permita a los usuarios ver el historial completo de transacciones relacionadas con su cuenta.
    Respuesta:
    ```json
    [
        {
            "id": 3,
            "amount": 1000.0,
            "transactionType": "CASH_TRANSFER",
            "transactionDate": 1729573542375,
            "sourceAccountNumber": "19b332",
            "targetAccountNumber": "ed9050"
        },
        {
            "id": 2,
            "amount": 20000.0,
            "transactionType": "CASH_WITHDRAWAL",
            "transactionDate": 1729573356164,
            "sourceAccountNumber": "19b332",
            "targetAccountNumber": "N/A"
        },
        {
            "id": 1,
            "amount": 100000.0,
            "transactionType": "CASH_DEPOSIT",
            "transactionDate": 1729573225112,
            "sourceAccountNumber": "19b332",
            "targetAccountNumber": "N/A"
        }
    ]
    ```
    Los tipos de transacciones que deben ser soportados por la aplicación son:
    - CASH_WITHDRAWAL
    - CASH_DEPOSIT
    - CASH_TRANSFER
    - SUBSCRIPTION
    - ASSET_PURCHASE
    - ASSET_SELL

    La fecha de la transacción también debe ser registrada en la base de datos.

#### Tarea 5: Operaciones de Mercado

Esta tarea se centra en las operaciones relacionadas con el mercado de valores, incluyendo la compra y venta de activos, la visualización de precios en tiempo real y el cálculo del valor neto.

Los endpoints que realizan acciones en las cuentas de los usuarios deben requerir el PIN y el JWT. Los endpoints que son meramente informativos serán públicos.

Para la compra y venta, se deben tomar los valores obtenidos en tiempo real utilizando la API proporcionada.

Para la compra y venta, es necesario realizar las operaciones necesarias para que los activos obtenidos con la cantidad invertida se mantengan en la cuenta.

- **Comprar Activos**: Implementa la funcionalidad para permitir a los usuarios comprar activos (por ejemplo, acciones) especificando el símbolo del activo, la cantidad a invertir y el PIN del usuario.
    Cuerpo de la solicitud:
    ```json
    {
        "assetSymbol": "GOLD",
        "pin": "1811",
        "amount": 1000.0
    }
    ```
    Respuesta:
    ```
    Compra de activo exitosa.
    ```
    Respuesta del endpoint informativo de activos del usuario:
    ```json
    {
    "GOLD": 0.6829947955796576
    }
    ```
    En caso de error o falta de fondos, se devuelve un código de estado 500 con el mensaje:
    ```
    Ocurrió un error interno al comprar el activo.
    ```

    El precio de compra también debe almacenarse, para que se pueda calcular una ganancia o pérdida en caso de una venta.

    También se debe enviar un correo electrónico con el asunto `Confirmación de Compra de Inversión` en el siguiente formato:

    ```
    Estimado Nuwe Test,

    Ha comprado con éxito 0.14 unidades de GOLD por un monto total de $50.00.

    Tenencias actuales de GOLD: 0.53 unidades

    Resumen de activos actuales:
    - GOLD: 0.53 unidades compradas a $1160.70

    Saldo de la cuenta: $63376.87
    Valor Neto: $63560.59

    Gracias por utilizar nuestros servicios de inversión.

    Atentamente,
    Equipo de Gestión de Inversiones
    ```

- **Vender Activos**: Permite a los usuarios vender activos de su cartera, especificando la cantidad a vender y verificando con el PIN del usuario.
    Cuerpo de la solicitud:
    ```json
    {
        "assetSymbol": "GOLD",
        "pin": "1811",
        "quantity": 0.3
    }
    ```
    Respuesta:
    ```
    Venta de activo exitosa.
    ```

    En caso de error o falta de fondos, se devuelve un código de estado 500 con el mensaje:
    ```
    Ocurrió un error interno al vender el activo.
    ```
    Respuesta del endpoint informativo de activos del usuario:
    ```json
    {
    "GOLD": 0.3829947955796576
    }
    ```

    Utilizando el precio de compra almacenado, se debe calcular el rendimiento de la transacción, ganancia o pérdida.

    También se debe enviar un correo electrónico con el asunto `Confirmación de Venta de Inversión` en el siguiente formato:

    ```
    Estimado Nuwe Test,

    Ha vendido con éxito 0.30 unidades de GOLD.

    Ganancia/Pérdida Total: $87.63

    Tenencias restantes de GOLD: 0.38 unidades

    Resumen de activos actuales:
    - GOLD: 0.38 unidades compradas a $1464.14

    Saldo de la cuenta: $78526.87
    Valor Neto: $79199.50

    Gracias por utilizar nuestros servicios de inversión.

    Atentamente,
    Equipo de Gestión de Inversiones
    ```

- **Valor Neto**: Proporciona a los usuarios una visión general de su valor neto combinando el saldo de efectivo y las tenencias de activos.
    Respuesta:
    ```
    79061.08163071838
    ```

- **Precios de Mercado en Tiempo Real**: Implementa endpoints para obtener los precios actuales del mercado para activos individuales y todo el mercado disponible. Esta información debe obtenerse de la API indicada.
    Respuesta para todos los activos:
    ```json
    {
        "AAPL": 81.05,
        "GOOGL": 1082.33,
        "TSLA": 75.71,
        "AMZN": 119.0,
        "MSFT": 161.23,
        "NFLX": 427.81,
        "FB": 11.68,
        "BTC": 8304.25,
        "ETH": 91.54,
        "XRP": 4.26,
        "GOLD": 1162.48,
        "SILVER": 4.24
    }
    ```
    Respuesta para activo individual:
    ```
    1175.95
    ```
    #### Tarea 6: Seguridad

    Esta tarea verifica la seguridad de la API comprobando el control de acceso para los endpoints públicos y privados.

    - **Endpoints Públicos**: Asegúrate de que los endpoints públicos como el inicio de sesión y el registro sean accesibles sin autenticación.
    - **Endpoints Privados Sin Autenticación**: Verifica que los endpoints privados devuelvan un error 401 o 403 si se accede sin autenticación. Debe mostrarse el mensaje "Acceso denegado".
    - **Endpoints Privados Con Autenticación**: Asegúrate de que los endpoints privados sean accesibles con un token JWT válido y realicen las acciones previstas.
    - **Seguridad de Contraseña**: La contraseña debe almacenarse cifrada usando BCrypt.

    #### Tarea 7: Manejo de Errores

    Esta tarea asegura que la aplicación maneje los errores de manera adecuada y proporcione retroalimentación apropiada al usuario.

    - **Correo Electrónico o Número de Teléfono Duplicado**: Asegúrate de que intentar registrar un usuario con un correo electrónico o número de teléfono existente resulte en un error 400 y un mensaje apropiado.
    - **Credenciales de Inicio de Sesión Inválidas**: Prueba que los intentos de inicio de sesión inválidos (por ejemplo, correo electrónico o contraseña incorrectos) devuelvan un estado 401 con el mensaje "Credenciales incorrectas".
    - **Validación de Contraseña**: Implementa reglas de validación de contraseña fuertes y devuelve mensajes de error específicos para las violaciones.
        Cuerpos de solicitud y respuestas:
        ```json
        {
            "name":"Nuwe Test",
            "password":"nuwetest1$",
            "email":"nuwe@nuwe.com",
            "address":"Main St",
            "phoneNumber":"666888115"
        }

        Respuesta: La contraseña debe contener al menos una letra mayúscula

        {
            "name":"Nuwe Test",
            "password":"Nuwetest",
            "email":"nuwe@nuwe.com",
            "address":"Main St",
            "phoneNumber":"666888115"
        }

        Respuesta: La contraseña debe contener al menos un dígito y un carácter especial

        {
            "name":"Nuwe Test",
            "password":"Nuwetest1",
            "email":"nuweeee@nuwe.com",
            "address":"Main St",
            "phoneNumber":"666888115"
        }

        Respuesta: La contraseña debe contener al menos un carácter especial

        {
            "name":"Nuwe Test",
            "password":"Nuwetest1 ",
            "email":"nuweeee@nuwe.com",
            "address":"Main St",
            "phoneNumber":"666888115"
        }

        Respuesta: La contraseña no puede contener espacios en blanco

        {
            "name":"Nuwe Test",
            "password":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa1$",
            "email":"nuweeee@nuwe.com",
            "address":"Main St",
            "phoneNumber":"666888115"
        }

        Respuesta: La contraseña debe tener menos de 128 caracteres

        {
            "name":"Nuwe Test",
            "password":"Test1$",
            "email":"nuweeee@nuwe.com",
            "address":"Main St",
            "phoneNumber":"666888115"
        }

        Respuesta: La contraseña debe tener al menos 8 caracteres
        ```

    - **Validación de Formato de Correo Electrónico**: Implementa reglas de validación de formato de correo electrónico fuertes.
        Cuerpo de solicitud:
        ```json
        {
            "name":"Nuwe Test",
            "password":"TestTest1$",
            "email":"nuwenuwe",
            "address":"Main St",
            "phoneNumber":"666888115"
        }
        ```
        Respuesta:
        ```
        Correo electrónico inválido: nuwenuwe
        ```
    - **Fondos insuficientes en la cuenta**: Cualquier transacción para la cual no haya fondos suficientes debe desencadenar el texto: **Saldo insuficiente** con un código de estado 400.

    ### Tarea 8: Suscripciones y Bot de Trading

    Esta tarea se centra en características avanzadas como la creación de suscripciones automáticas y la habilitación de un bot de trading para manejar inversiones.

    - **Crear Suscripción**: Implementa una función que permita a los usuarios suscribirse a pagos periódicos en un intervalo establecido. En este caso, el intervalo será en segundos, para poder comprobar el correcto funcionamiento de la aplicación.
        Cuerpo de solicitud:
        ```json
        {
            "pin": "1811",
            "amount":"100",
            "intervalSeconds": 5
        }
        ```
        Respuesta:
        ```
        Suscripción creada exitosamente.
        ```
        El correcto funcionamiento de las suscripciones depende de si se crea correctamente, se ejecuta en el intervalo indicado simulando una suscripción o domiciliación y la cantidad de dinero en la cuenta disminuye hasta que no quede nada.

        Estas transacciones también deben guardarse como otras transacciones con el tipo de transacción apropiado mencionado anteriormente.

    - **Bot de Inversión Automática**: Permite a los usuarios activar un bot de inversión automática que compra o vende activos automáticamente según las condiciones del mercado (por ejemplo, fluctuaciones de precios).
        Cuerpo de solicitud:
        ```json
        {
            "pin": "1811"
        }
        ```
        Respuesta:
        ```
        Inversión automática habilitada exitosamente.
        ```

        En este caso, se deben crear reglas para comprar o vender activos automáticamente según las fluctuaciones del mercado, es decir, los precios devueltos por la API en tiempo real. Se recomienda utilizar reglas con pequeñas variaciones para probar el correcto funcionamiento.

        Debes usar un intervalo de tiempo que no comprometa el rendimiento de la aplicación, por ejemplo, 30 segundos. Es decir, cada 30 segundos el bot verifica si ha habido alguna fluctuación en el mercado para los activos que posee el usuario. Si es así, compra o vende dependiendo de si el activo baja o sube.

        Caso de estudio:

        El usuario tiene ORO y lo ha comprado a 1000. Si el precio baja a 800, compra una pequeña cantidad ya que podría apreciarse en valor. Si el precio sube a 1200, vende una parte de los activos para obtener rentabilidad.

        Esta rentabilidad debe ser calculada.













