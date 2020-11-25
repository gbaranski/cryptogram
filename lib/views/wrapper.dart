import 'package:cryptogram/models/stellar.dart';
import 'package:cryptogram/views/chat/index.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'initialize/index.dart';

class Wrapper extends StatelessWidget {
  static const route = '/';

  @override
  Widget build(BuildContext context) {
    final Stellar stellarModel = Provider.of<Stellar>(context);
    if (stellarModel.keyPair == null) {
      Navigator.of(context).pushNamed(InitializeView.route);
    }
    return ChatView();
  }
}
