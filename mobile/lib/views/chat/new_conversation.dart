import 'package:cryptogram/models/account.dart';
import 'package:cryptogram/services/blockchain.dart';
import 'package:cryptogram/services/crypto.dart';
import 'package:cryptogram/views/account/get_password.dart';
import 'package:flutter/material.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Row, Account;

class NewConversationView extends StatefulWidget {
  static const route = '/new_conversation';
  final Account account;
  NewConversationView(this.account);

  @override
  _NewConversationViewState createState() => _NewConversationViewState();
}

class _NewConversationViewState extends State<NewConversationView> {
  final _formKey = GlobalKey<FormState>();

  String destinationAccountID;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: Container(
        child: Form(
          key: _formKey,
          autovalidateMode: AutovalidateMode.onUserInteraction,
          child: ListTile(
            trailing: IconButton(
              onPressed: () async {
                if (!_formKey.currentState.validate()) return;
                _formKey.currentState.save();
                final userPassword = await getUserPassword(context);
                final decryptedSecretSeed = await Crypto.decryptSecretSeed(
                    widget.account.secretSeed, userPassword);
                print("Decrypted secret seed: $decryptedSecretSeed");
                final senderKeyPair =
                    KeyPair.fromSecretSeed(decryptedSecretSeed);

                final result = await Blockchain.sendMessage(
                    destination: destinationAccountID,
                    message: "Hello world",
                    senderKeypair: senderKeyPair);

                print("Success: ${result.success}");
                print("Sender: ${widget.account.accountID}");
                print("Destination: $destinationAccountID");
              },
              icon: Icon(Icons.send),
            ),
            title: TextFormField(
                decoration: const InputDecoration(labelText: 'Account ID'),
                autocorrect: false,
                onSaved: (String value) {
                  setState(() {
                    destinationAccountID = value;
                  });
                },
                validator: (value) {
                  try {
                    if (value == null)
                      throw new Exception('AccountID cannot be empty');
                    KeyPair.fromAccountId(value);
                    return null;
                  } catch (e) {
                    return "Invalid: ${e.message.toLowerCase()}";
                  }
                }),
          ),
        ),
      ),
    );
  }
}
