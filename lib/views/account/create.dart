import 'package:flutter/material.dart';

class CreateAccountView extends StatelessWidget {
  static const String route = '/account/create';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: Text("Create new wallet"),
    );
  }
}
