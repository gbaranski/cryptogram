import 'package:cryptogram/route_generator.dart';
import 'package:cryptogram/services/database.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // Temporary solution
  await DatabaseService.init();
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Cryptogram',
      initialRoute: '/',
      onGenerateRoute: RouteGenerator.generateRoute,
      themeMode: ThemeMode.dark,
      theme: ThemeData(
          brightness: Brightness.dark,
          colorScheme: const ColorScheme.dark(),
          bottomNavigationBarTheme: BottomNavigationBarThemeData(
              selectedItemColor: Colors.indigoAccent),
          floatingActionButtonTheme: FloatingActionButtonThemeData(
              backgroundColor: Colors.deepPurpleAccent,
              foregroundColor: Colors.white.withAlpha(240)),
          scaffoldBackgroundColor: Colors.grey.shade900,
          textTheme: GoogleFonts.montserratTextTheme(
              Theme.of(context).primaryTextTheme)),
    );
  }
}
