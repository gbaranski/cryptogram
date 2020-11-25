import 'package:cryptogram/models/stellar.dart';
import 'package:cryptogram/views/account/create.dart';
import 'package:cryptogram/views/account/index.dart';
import 'package:cryptogram/views/wrapper.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';

import 'views/chat/index.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Cryptogram',
      initialRoute: '/',
      routes: {
        // '/'
        Wrapper.route: (context) => Wrapper(),
        // '/account'
        AccountView.route: (context) => AccountView(),
        CreateAccountView.route: (context) => CreateAccountView(),
        // '/chat'
        ChatView.route: (context) => ChatView(),
      },
      theme: ThemeData(
          primarySwatch: Colors.indigo,
          textTheme:
              GoogleFonts.openSansTextTheme(Theme.of(context).textTheme)),
      home: ChangeNotifierProvider(
        create: (_) => new Stellar(),
        child: Wrapper(),
      ),
    );
  }
}
