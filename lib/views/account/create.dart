import 'package:cryptogram/models/account.dart';
import 'package:cryptogram/services/crypto.dart';
import 'package:cryptogram/services/database.dart';
import 'package:cryptogram/views/account/get_password.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Row, Account;

class CreateAccountView extends StatefulWidget {
  static const String route = '/account/create';

  @override
  _CreateAccountViewState createState() => _CreateAccountViewState();
}

class _CreateAccountViewState extends State<CreateAccountView> {
  KeyPair _keyPair = KeyPair.random();

  Future<void> addAccount(BuildContext context, String customName) async {
    final password = await getUserPassword(context);
    if (password == null) return;
    final encryptedSecret =
        await Crypto.encryptSecretSeed(_keyPair.secretSeed, password);
    final account = Account(
        accountID: _keyPair.accountId,
        customName: customName,
        secretSeed: encryptedSecret);

    DatabaseService.addAccount(account);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: Container(
        margin: EdgeInsets.symmetric(horizontal: 20, vertical: 20),
        child: Column(
          children: [
            ListTile(
              leading: Icon(MdiIcons.accountKey),
              onTap: () {
                Clipboard.setData(ClipboardData(text: _keyPair.accountId));
                HapticFeedback.lightImpact();
                ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
                  content: Text("Copied AccountID to clipboard"),
                ));
              },
              title: Text("Account ID"),
              subtitle: Text(_keyPair.accountId),
            ),
            ListTile(
              leading: Icon(MdiIcons.keyStar),
              onTap: () {
                Clipboard.setData(ClipboardData(text: _keyPair.accountId));
                HapticFeedback.lightImpact();
                ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
                  content: Text(
                      "Copied secret to clipboard, save it in secure place"),
                ));
              },
              title: Text("Secret seed"),
              subtitle: Text(_keyPair.secretSeed),
            ),
            ButtonBar(
              alignment: MainAxisAlignment.spaceEvenly,
              children: [
                OutlinedButton(
                  onPressed: () {
                    setState(() {
                      _keyPair = KeyPair.random();
                    });
                  },
                  child: Text("Generate new key pair"),
                ),
                ElevatedButton(
                  onPressed: () async {
                    addAccount(context, "SomeCustomName");
                  },
                  child: Text("Continue"),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }
}
