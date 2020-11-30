import 'package:cryptogram/models/account.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Account;

abstract class Blockchain {
  static final _sdk = StellarSDK.TESTNET;

  static Future<bool> fundAccount(Account account) =>
      FriendBot.fundTestAccount(account.accountID);

  static Future<AccountResponse> getAccountInfo(Account account) =>
      _sdk.accounts.account(account.accountID);
}
