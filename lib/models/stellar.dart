import 'package:flutter/material.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart';

class Stellar extends ChangeNotifier {
  final StellarSDK _stellarSDK = StellarSDK.TESTNET;
  KeyPair keyPair;

  KeyPair generateKeyPair() {
    keyPair = KeyPair.random();
    notifyListeners();
    return keyPair;
  }

  Future<bool> fundAccount() => FriendBot.fundTestAccount(keyPair.accountId);
}
