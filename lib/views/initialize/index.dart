import 'package:cryptogram/views/account/create.dart';
import 'package:flutter/material.dart';

class InitializeView extends StatelessWidget {
  static const route = '/initialize';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        margin: EdgeInsets.symmetric(horizontal: 50, vertical: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.end,
          children: [
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pushNamed(CreateAccountView.route);
              },
              child: Text("Create a new account"),
            ),
            TextButton(
              onPressed: () {},
              child: Text("I already have account"),
            )
          ],
        ),
      ),
    );
  }
}
