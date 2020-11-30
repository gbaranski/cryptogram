import 'package:cryptogram/models/account.dart';
import 'package:cryptogram/views/account/accounts_list.dart';
import 'package:cryptogram/views/chat/new_conversation.dart';
import 'package:flutter/material.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Account;

enum MiscAction {
  switch_account,
}

class ChatView extends StatefulWidget {
  ChatView(this.account);
  final Account account;

  @override
  _ChatViewState createState() => _ChatViewState();
}

class _ChatViewState extends State<ChatView> {
  AccountResponse accountData;

  void onMiscAction(BuildContext context, MiscAction action) {
    if (action == MiscAction.switch_account) {
      Navigator.pushReplacementNamed(context, AccountsList.route);
    }
  }

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () =>
            Navigator.pushNamed(context, NewConversationView.route),
        icon: Icon(MdiIcons.messagePlusOutline),
        label: Text("Start chat"),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await Future.delayed(Duration(milliseconds: 500));
        },
        child: CustomScrollView(
          slivers: [
            SliverSafeArea(
                minimum: EdgeInsets.symmetric(vertical: 30),
                sliver: SliverAppBar(
                  flexibleSpace: Container(
                    margin: EdgeInsets.symmetric(horizontal: 10),
                    child: Ink(
                      decoration: BoxDecoration(
                          borderRadius: BorderRadius.circular(15),
                          color: Color(0xFF313131)),
                      child: ListTile(
                        onTap: () {},
                        trailing: PopupMenuButton<MiscAction>(
                          padding: EdgeInsets.zero,
                          icon: Icon(Icons.more_vert),
                          onSelected: (action) => onMiscAction(context, action),
                          itemBuilder: (context) => [
                            PopupMenuItem(
                                value: MiscAction.switch_account,
                                child: ListTile(
                                  visualDensity: VisualDensity.compact,
                                  contentPadding: EdgeInsets.zero,
                                  leading: Icon(MdiIcons.accountSwitch),
                                  title: Text("Switch account"),
                                )),
                          ],
                        ),
                      ),
                    ),
                  ),
                )),
          ],
        ),
      ),
    );
  }
}
