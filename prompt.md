refazer o processo de notificação que está bem confuso.
Necessidades do novo sistema:
1) Existem 3 tipos possíveis de notificação:
    a) e-mail - utilizando o email_adapater.go
    b) sms - utilizando o sms adapter.go
    c) pushNotification ou fcm - utilizando o fcm_adapter.go
2) como chamar a notificador
    a) type - define que tipo de notificação será enviada (sms/fcm/email)
    a) from - opcional e será usada na notificação de email
    b) to - obrigatória para email e sms. contem o numero de telefone para sms, o endereço de email para email
    c) subject - obrigatória para e-mail. contem o subject do email, title do fcm
    d) body - pbrigatorio para todos. contem o corpo da mensagem
    e) imageUrl - opcional e conterá imageUrl do fcm
    f) token - necessário para fcm, conterá o deviceToken
4) Construa este novo sistema em substituição ao notificaton_handler e notification_sender.
5) Verifique as rotinas que chamavam os notificaton_handler e/ou notification_sender e altere para chamar esta novo rotina