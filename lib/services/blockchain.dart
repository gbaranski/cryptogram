import 'package:cryptogram/models/account.dart';
import 'package:flutter/material.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Account;
import 'package:url_launcher/url_launcher.dart';

abstract class Blockchain {
  static final _sdk = StellarSDK.TESTNET;

  static Future<bool> fundAccount(Account account) =>
      FriendBot.fundTestAccount(account.accountID);

  static Future<AccountResponse> getAccountInfo(Account account) =>
      _sdk.accounts.account(account.accountID);

  static Future<void> openAccountInfo(Account account) async {
    final String url =
        "http://testnet.stellarchain.io/address/${account.accountID}";
    if (await canLaunch(url)) {
      await launch(url);
    } else {
      throw new Exception('Could not open URL');
    }
  }

  static Future<SubmitTransactionResponse> sendMessage(
      {@required KeyPair senderKeypair,
      @required String destination,
      @required String message}) async {
    AccountResponse sender =
        await _sdk.accounts.account(senderKeypair.accountId);
    final memo = MemoText(message);
    Transaction transaction = new TransactionBuilder(sender)
        .addOperation(
            PaymentOperationBuilder(destination, Asset.NATIVE, "100").build())
        .addMemo(memo)
        .build();
    transaction.sign(senderKeypair, Network.TESTNET);

    return _sdk.submitTransaction(transaction);
  }

  static Stream<OperationResponse> paymentsStream(Account account) {
    return _sdk.payments.forAccount(account.accountID).stream();
  }
}
