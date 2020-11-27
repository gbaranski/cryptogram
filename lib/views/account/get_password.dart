import 'package:flutter/material.dart';

class GetPasswordDialog extends StatefulWidget {
  @override
  _GetPasswordDialogState createState() => _GetPasswordDialogState();
}

class _GetPasswordDialogState extends State<GetPasswordDialog> {
  final _formKey = GlobalKey<FormState>();
  String _password = "";

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text("Enter password for encryption"),
      content: Form(
        key: _formKey,
        autovalidateMode: AutovalidateMode.onUserInteraction,
        child: TextFormField(
          decoration: const InputDecoration(labelText: "Password"),
          obscureText: true,
          autocorrect: false,
          enableSuggestions: false,
          onSaved: (value) => setState(() => _password = value),
          validator: (value) {
            if (value == null) return "Password cannot be empty";
            if (value.length < 8)
              return "Password must be longer or equal to 8";
            return null;
          },
        ),
      ),
      actions: [
        FlatButton(
          child: Text("Submit"),
          onPressed: () {
            if (!_formKey.currentState.validate()) return;
            _formKey.currentState.save();
            Navigator.pop(context, _password);
          },
        )
      ],
    );
  }
}

Future<String> getUserPassword(BuildContext context) async {
  return await showDialog<String>(
      context: context,
      builder: (BuildContext context) => GetPasswordDialog(),
      barrierDismissible: false);
}
