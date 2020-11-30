import 'dart:convert';
import 'dart:typed_data';

import 'package:cryptography/cryptography.dart';

abstract class Crypto {
  static const secretSeedCipher = CipherWithAppendedMac(aesCtr, Hmac(sha256));
  static final nonce = Nonce([6, 8, 1, 5, 0, 4, 5, 8, 7, 9, 1, 2]);

  static Future<Uint8List> encryptSecretSeed(
      String secretSeed, String password) async {
    final secretKey = await sha256.hash(utf8.encode(password));
    return secretSeedCipher.encrypt(utf8.encode(secretSeed),
        secretKey: SecretKey(secretKey.bytes), nonce: nonce);
  }
}
