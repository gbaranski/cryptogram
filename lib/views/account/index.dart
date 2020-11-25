import 'package:flutter/material.dart';

class AccountView extends StatelessWidget {
  static const route = '/account';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        margin: EdgeInsets.symmetric(horizontal: 50, vertical: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.end,
          children: [Text("Account view")],
        ),
      ),
    );
  }
}
