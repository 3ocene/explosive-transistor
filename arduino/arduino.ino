void setup(void) {
  Serial.begin(9600);
  for (int i = 0; i <= 13; i++) {
    pinMode(i, OUTPUT);
  }
}

int pinNumber = 0;
int value;

void loop() {
  switch (Serial.read()) {
    case 0:
      waitForData();
      Serial.write(digitalRead(Serial.read()));
      break;
    case 1:
      waitForData();
      pinNumber = Serial.read();
      waitForData();
      value = Serial.read();
      digitalWrite(pinNumber, toOutput(value));
      break;
    case -1:
      return;
  }
}

void waitForData() {
  while (!Serial.available()) {
    ;
  }
}

int toOutput(int v) {
  if (v == 0) {
    return LOW;
  }
  if (v == 1) {
    return HIGH;
  }
}

