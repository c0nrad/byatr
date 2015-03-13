# ECB Byte-At-A-Time Decryption

This article about a crypto attack called ECB byte-at-a-time decryption. We will cover the scenario where this attack works, necessary background infromation about the attack, some points to consider, and then go over a tool I'm releasing called b44tr.

## ECB

There are a bunch of different ways to encrypt data. But most of the methods fall under two broad categories: asynmetric encryption and symetric encryption. Asynmetric ciphers require two different keys, one for encryption and one for decryption. This allows you to do some pretty neat stuff. For example you could gerenate a key pair, then give everyone in the world your key used for encryption. Everyone in the world can then securely send you information, and since you're the only person that has the key used for decryption, only you will be able to see what all the messages say. This is also how you can securely send websites like amazon.com your credit card information without every person in the world also knowing what it is. But I digress.

We'll be talking about a type of symmetric key encryption called block ciphers. In symmetric key encryption, you use the same key for encrypting and decrypting data.
