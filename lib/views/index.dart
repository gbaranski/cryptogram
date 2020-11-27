import 'package:flutter/material.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart' hide Row;

enum AccountAction {
  delete,
  edit,
}

class IndexView extends StatefulWidget {
  static const route = '/';

  @override
  _IndexViewState createState() => _IndexViewState();
}

class _IndexViewState extends State<IndexView> {
  List<Account> _accounts; // to be implemented

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
        child: Icon(MdiIcons.accountPlus),
      ),
      body: CustomScrollView(
        physics:
            AlwaysScrollableScrollPhysics().applyTo(BouncingScrollPhysics()),
        slivers: [
          SliverAppBar(
            title: Text("Pick account"),
          ),
          SliverList(
            delegate: SliverChildBuilderDelegate(
                (context, i) => ListTile(
                      onTap: () {},
                      leading: Icon(MdiIcons.accountKey),
                      trailing: PopupMenuButton<AccountAction>(
                        icon: Icon(Icons.more_vert),
                        itemBuilder: (context) => [
                          PopupMenuItem(
                              value: AccountAction.delete,
                              child: ListTile(
                                visualDensity: VisualDensity.compact,
                                contentPadding: EdgeInsets.zero,
                                leading: Icon(Icons.delete),
                                title: Text("Delete"),
                              )),
                          PopupMenuItem(
                              value: AccountAction.edit,
                              enabled: false,
                              child: ListTile(
                                visualDensity: VisualDensity.compact,
                                contentPadding: EdgeInsets.zero,
                                leading: Icon(Icons.edit),
                                title: Text(
                                  "Edit",
                                  style: TextStyle(color: Colors.black38),
                                ),
                              )),
                        ],
                      ),
                      title: Text(
                        KeyPair.random().accountId,
                        overflow: TextOverflow.ellipsis,
                      ),
                      subtitle: Text("Grzegorz"),
                    ),
                childCount: 2),
          )
        ],
      ),
    );
  }
}
