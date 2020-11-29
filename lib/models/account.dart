import 'package:flutter/material.dart';

class Account {
  final String accountID;
  final String customName;
  final String secretSeed;

  Map<String, dynamic> toMap() {
    return {
      'accountID': accountID,
      'customName': customName,
      'secretSeed': secretSeed
    };
  }

  factory Account.fromMap(Map<String, dynamic> map) {
    return Account(
      accountID: map['accountID'],
      customName: map['customName'],
      secretSeed: map['secretSeed'],
    );
  }

  Account({@required this.accountID, @required this.customName, @required this.secretSeed});
}
