// Copyright (c) 2017, raz. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

import 'dart:html';
import 'package:intl/intl.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';
//import 'dart:async';

var url = 'http://127.0.0.1:8000';
var urlInsertForm = '$url/InsertNewGasSnapForm';
var urlSvgMem = '$url/GetSvgFromMemory';
var urlSvgSaved = "$url/GetSvg";
var urlUpdateGas = '$url/UpdateGas';
//var urlSvgFile = "$url/GetSvg";
var urlInsertFormString = '$url/InsertNewGasSnapFormString';

BrowserClient client = new BrowserClient();
DivElement DivOutput;

InputElement IDate;
InputElement IcPrice;
InputElement IcMeters;
InputElement IcMetersToBeUpdated;
InputElement IcMetersNewVal;

ImageElement ImageGraph1;
ButtonElement BtnSubmit;
ButtonElement BtnInsert;
ButtonElement BtnDeleteLast;
ButtonElement BtnUpdate;
AnchorElement linkSvgFile;

main() async {
  DivOutput = querySelector('#output');
  IDate = querySelector('#date');

  IcPrice = querySelector("#cPrice");
  document.onReadyStateChange.listen(last_cPrice);

  ImageGraph1 = querySelector('#graph1');

  IcMeters = querySelector("#cMeters");

  BtnInsert = querySelector('#insert');
  BtnInsert.onClick.listen(insertGas);

  BtnSubmit = querySelector('#submit');
  BtnSubmit.onClick.listen(saveToFiles);

  IcMetersToBeUpdated = querySelector("#cMetersToBeUpdated");
  IcMetersNewVal = querySelector('#cMetersNewVal');

  BtnUpdate = querySelector('#bUpdate');
  BtnUpdate.onClick.listen(updateGas);

  BtnDeleteLast = querySelector('#deleteLast');
  BtnDeleteLast.onClick.listen(deleteLast);

  IDate.value = setDate();

  linkSvgFile = querySelector("#linkSvgFile");
  linkSvgFile.onClick.listen(linkToSavedSvg);
}

int nr = 0;

linkToSavedSvg(_) {
  ImageGraph1.src = "$urlSvgSaved#$nr";
  nr++;
}

saveToFiles(_) async {
  Response resp = await client.get('$url/SaveToFiles');
  if (resp.statusCode != 200) {
    DivOutput.text = resp.body;
    return;
  }

  DivOutput.text = "Saved to files.";
}

refreshSvg() {
  ImageGraph1.src = "$urlSvgMem#$nr";
  nr++;
}

deleteLast(_) async {
  Response resp;

  try {
    resp = await client.get("$url/DeleteLastGas");
  } catch (e) {
    DivOutput.text = e.toString();
  }

  if (resp.body != 'ok') {
    DivOutput.text = "Could not delete err: ${resp.body}";
    return;
  }

  DivOutput.text += "Deleted last entry";
  refreshSvg();
}

last_cPrice(_) async {
  Response resp = await client.get("$url/GetLast_cPrice");

  Element IcPriceParent = IcPrice.parent;
  IcPriceParent.classes.add('is-dirty');

  IcPrice.value = resp.body;
}

insertGas(_) async {
  var response = await client
      .post(urlInsertFormString,
          body: "${IDate.value}|${IcMeters.value}|${IcPrice.value}")
      .catchError((err) {
    print(err);
    DivOutput.text = err;
  });

  if (response.statusCode != 200) {
    print('err rbody: ${response.body}');
    DivOutput.text = response.body;
    return;
  }

//  HttpRequest respPostForm = await HttpRequest.postFormData(urlInsertForm, {
//    "date": IDate.value,
//    "cMeters": IcMeters.value,
//    "cPrice": IcPrice.value
//  });
//
//  if (respPostForm.responseText != 'ok') {
//    DivOutput.text = respPostForm.responseText;
//    return;
//  }

//  Response respSave = await client.get("$url/SaveToFiles");
//
//  if (respSave.body != "ok") {
//    DivOutput.text = respSave.body;
//    return;
//  }
//
//  DivOutput.text = "Saved to file.";

  DivOutput.text = 'inserted ${IDate.value} ${IcMeters.value} ${IcPrice.value}';
  refreshSvg();
}

updateGas(_) async {
  try {
    Response resp = await client.post(urlUpdateGas,
        body: "${IcMetersToBeUpdated.value} ${IcMetersNewVal.value}");

    if (resp.statusCode != 200) {
      DivOutput.text = "updateGas: ${resp.body}";
    } else {
      DivOutput.text = "Updated";
      refreshSvg();
    }
  } catch (e) {
    DivOutput.text = e.toString();
  }
}

//Future<String> sleep1() {
//  return new Future.delayed(const Duration(seconds: 2),()=>'');
//}

String setDate() {
  var now = new DateTime.now();
  var formatter = new DateFormat('yyyy/MM/dd HH:mm');
  return formatter.format(now);
}
