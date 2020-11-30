import 'package:flutter/material.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Row;

class NewConversationView extends StatefulWidget {
  static const route = '/new_conversation';
  @override
  _NewConversationViewState createState() => _NewConversationViewState();
}

class _NewConversationViewState extends State<NewConversationView> {
  final _formKey = GlobalKey<FormState>();

  String accountID;

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
              onPressed: () {},
              icon: Icon(Icons.send),
            ),
            title: TextFormField(
                decoration: const InputDecoration(labelText: 'Account ID'),
                autocorrect: false,
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

            //   IconButton(
            //     onPressed: () {},
            //     icon: Icon(Icons.arrow_right),
            //   )
          ),
        ),
      ),
    );
  }
}
