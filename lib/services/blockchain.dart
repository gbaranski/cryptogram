import 'package:cryptogram/models/account.dart';
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
}
