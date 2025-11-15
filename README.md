# KON-KHMER-GENZ

## 🇬🇧 Official English Version

## Server Resilience Tester (Official README)

**KON-KHMER-GENZ** is a project created for **educational**, **research**, and **server resilience testing** purposes. It is designed to help administrators evaluate how well their systems handle high traffic loads.

⚠️ **IMPORTANT NOTICE**

* This tool **must not** be used to attack or disrupt any system, network, or service that you do **not** own or do **not** have clear permission to test.
* Misuse of stress-testing or high‑traffic simulation tools can be **illegal** and may lead to serious consequences.
* The developer and contributors are **not responsible** for any misuse of this tool.

## Purpose of This Project

* To help server owners test stability under heavy load.
* To study performance behavior in controlled environments.
* To assist students and researchers in learning about traffic simulation and resilience testing.

## Features (Safe Overview)

* Supports different testing modes:

  * **TCP Connection Load Test** (Mode 1)
  * **UDP High-Throughput Test** (Mode 2)
  * **HTTP Load Simulation** (Mode 3)
* Customizable parameters such as connections, duration, timeout, and packet size.
* Optional proxy file usage for authorized HTTP load simulations.

## Usage (Safe & Official)

This tool accepts the following arguments:

```
<host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]
```

If the argument count is incorrect, the tool displays:

```
Usage: <program> <host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]
---

## 🇰🇭 កំណែជាភាសាខ្មែរ (Khmer Version)

### **KON-KHMER-GENZ**

**KON-KHMER-GENZ** គឺជាគម្រោងមួយសម្រាប់ **ការសិក្សា**, **ការស្រាវជ្រាវ**, និង **ការធ្វើតេស្តសមត្ថភាពServer** នៅក្រោមបន្ទុកខ្ពស់។ វាត្រូវបានបង្កើតឡើងសម្រាប់ម្ចាស់ម៉ាស៊ីនឬអ្នកអភិវឌ្ឍន៍ប្រព័ន្ធ ដើម្បីត្រួតពិនិត្យការអនុវត្តន៍ និងការអាំបានរបស់ម៉ាស៊ីនពេលមានចរាចរណ៍ច្រើន។

⚠️ **សារាជាពិពណ៌នា**
- មិនត្រូវប្រើឧបករណ៍នេះ ដើម្បីវាយប្រហារ ឬបង្កការរំខានលើប្រព័ន្ធណាមួយ **ដែលអ្នកមិនមែនជាម្ចាស់** ឬ **មិនមានការអនុញ្ញាត** ណាមួយឡើយ។
- ការប្រើឧបករណ៍សម្រាប់បង្កើតចរាចរណ៍ខ្ពស់ (Load/Stress Test) លើគេហទំព័រ ឬម៉ាស៊ីនរបស់អ្នកដទៃ អាចជាអំពើ **ខុសច្បាប់** និងមានទោសធ្ងន់ធ្ងរ។
- អ្នកបង្កើត និងអ្នកចូលរួមក្នុងគម្រោងនេះ **មិនទទួលខុសត្រូវ** ចំពោះការប្រើប្រាស់ខុសគោលបំណងឡើយ។

### **គោលបំណងរបស់គម្រោង**
- ប្រើសម្រាប់តេស្តកម្លាំងម៉ាស៊ីននៅក្រោមបន្ទុកខ្ពស់។
- សម្រាប់សិក្សាពីការប្រតិបត្តិរបស់Server នៅពេលមានចរាចរណ៍ច្រើន។
- ជួយសិស្ស-និស្សិត និងអ្នកស្រាវជ្រាវក្នុងការសិក្សាអំពីSimulation Traffic។

### **មុខងារ**
- គាំទ្រមូដធ្វើតេស្ត៖
  - TCP Connection Load Test (Mode 1)
  - UDP High-Throughput Test (Mode 2)
  - HTTP Load Test (Mode 3)
- កំណត់Connections, Timeout, Packetsize និងProxy File បាន។

### **ការប្រើប្រាស់ (Usage)**
```

<host> <port> <mode> <connections> <seconds> <timeout> <packetsize> [proxyfile.txt]

```

### **Mode List**
```

[1] TCP-Connection Load Test
[2] UDP Test (High Throughput)
[3] HTTP Load Test (High RPS)

```

### **ឧទាហរណ៍ (Examples)**
**TCP Example:**
```

program.exe 127.0.0.1 8080 1 50 60 3 1024

```

**UDP Example:**
```

program.exe 127.0.0.1 8080 2 100 45 3 2048

```

**HTTP Example (មានProxy File):**
```

program.exe 127.0.0.1 80 3 200 30 5 0 proxies.txt

```

---

ប្រសិនបើអ្នកចង់បន្ថែម **Logo**, **Version History**, **Installation Guide**, ឬ **Banner Khmer Style**, ប្រាប់ខ្ញុំបានគ្រប់ពេល!-----------------------------------------------
|                M O D E   L I S T               |
|      [1] TCP-Connection Load Test              |
|      [2] UDP Test (High Throughput)            |
|      [3] HTTP Load Test (High RPS)             |
--------------------------------------------------
```

### Example Commands (For Authorized Testing Only)

**TCP Mode Example:**

```
program.exe 127.0.0.1 8080 1 50 60 3 1024
```

**UDP Mode Example:**

```
program.exe 127.0.0.1 8080 2 100 45 3 2048
```

**HTTP Mode Example (With Proxy File):**

```
program.exe 127.0.0.1 80 3 200 30 5 0 proxies.txt
```

## Recommended Safe Usage

* Only test on servers you control or have written permission to test.
* Use within closed networks or sandbox environments.
* Always follow your local laws and ethical guidelines.

## Screenshot

*(Screenshot: 2025-11-15 18:32:05)*

## Official Note

This README is created to provide an **official, safe, and responsible** description of the KON-KHMER-GENZ project. The project should be used strictly within legal and ethical boundaries.

---

If you want to update sections like **installation**, **contributors**, **version history**, or **UI previews**, just let me know!
