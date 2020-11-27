import 'package:flutter/material.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class RestoreAccountView extends StatelessWidget {
  static const String route = '/account/restore';

  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: Container(
        margin: EdgeInsets.symmetric(horizontal: 30),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                decoration: const InputDecoration(
                    icon: Icon(MdiIcons.key), labelText: "Secred seed"),
                autocorrect: false,
                validator: (value) {
                  try {
                    KeyPair.fromSecretSeed(value);
                    return null;
                  } catch (e) {
                    return "Invalid: ${e.message.toLowerCase()}";
                  }
                },
              ),
              SizedBox(
                height: 20,
              ),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    if (!_formKey.currentState.validate()) return;
                  },
                  child: const Text("Restore account"),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
