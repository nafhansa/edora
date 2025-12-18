import 'dart:convert';

class SignalPoint {
  final int freq;
  final double impedance;
  final double? phase;

  SignalPoint({required this.freq, required this.impedance, this.phase});

  factory SignalPoint.fromJson(Map<String, dynamic> j) => SignalPoint(
        freq: j['freq'] is int ? j['freq'] : (j['freq'] as num).toInt(),
        impedance: (j['impedance'] as num).toDouble(),
        phase: j.containsKey('phase') ? (j['phase'] as num).toDouble() : null,
      );

  Map<String, dynamic> toJson() => {
        'freq': freq,
        'impedance': impedance,
        if (phase != null) 'phase': phase,
      };
}

class AnalysisResult {
  final double? bmd;
  final double? tScore;
  final String? classification;

  AnalysisResult({this.bmd, this.tScore, this.classification});

  factory AnalysisResult.fromJson(Map<String, dynamic> j) => AnalysisResult(
        bmd: j['bmd'] == null ? null : (j['bmd'] as num).toDouble(),
        tScore: j['t_score'] == null ? null : (j['t_score'] as num).toDouble(),
        classification: j['class'] as String?,
      );

  Map<String, dynamic> toJson() => {
        if (bmd != null) 'bmd': bmd,
        if (tScore != null) 't_score': tScore,
        if (classification != null) 'class': classification,
      };
}

class Location {
  final double? lat;
  final double? lng;

  Location({this.lat, this.lng});

  factory Location.fromJson(Map<String, dynamic> j) => Location(
        lat: j['lat'] == null ? null : (j['lat'] as num).toDouble(),
        lng: j['lng'] == null ? null : (j['lng'] as num).toDouble(),
      );

  Map<String, dynamic> toJson() => {
        if (lat != null) 'lat': lat,
        if (lng != null) 'lng': lng,
      };
}

class Reading {
  final String deviceSerial;
  final String patientId;
  final String doctorId;
  final DateTime timestamp;
  final Location? location;
  // Raw signal data can be a list of SignalPoint or arbitrary objects; keep flexible
  final List<dynamic> rawSignalData;
  final AnalysisResult? analysis;

  Reading({
    required this.deviceSerial,
    required this.patientId,
    required this.doctorId,
    required this.timestamp,
    this.location,
    required this.rawSignalData,
    this.analysis,
  });

  factory Reading.fromJson(Map<String, dynamic> j) {
    final readings = <dynamic>[];
    if (j['readings'] is List) {
      for (final r in (j['readings'] as List)) {
        if (r is Map<String, dynamic>) {
          readings.add(SignalPoint.fromJson(r));
        } else if (r is Map) {
          readings.add(SignalPoint.fromJson(Map<String, dynamic>.from(r)));
        } else {
          readings.add(r);
        }
      }
    } else if (j['raw_signal_data'] is List) {
      for (final r in (j['raw_signal_data'] as List)) {
        if (r is Map<String, dynamic>) {
          readings.add(SignalPoint.fromJson(r));
        } else if (r is Map) {
          readings.add(SignalPoint.fromJson(Map<String, dynamic>.from(r)));
        } else {
          readings.add(r);
        }
      }
    }

    return Reading(
      deviceSerial: j['device_serial'] as String? ?? j['deviceSerial'] as String? ?? '',
      patientId: j['patient_id'] as String? ?? j['patientId'] as String? ?? '',
      doctorId: j['doctor_id'] as String? ?? j['doctorId'] as String? ?? '',
      timestamp: j['timestamp'] == null ? DateTime.now() : DateTime.parse(j['timestamp'] as String),
      location: j['location'] == null ? null : Location.fromJson(Map<String, dynamic>.from(j['location'])),
      rawSignalData: readings,
      analysis: j['analysis'] == null ? null : AnalysisResult.fromJson(Map<String, dynamic>.from(j['analysis'])),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'device_serial': deviceSerial,
      'patient_id': patientId,
      'doctor_id': doctorId,
      'timestamp': timestamp.toUtc().toIso8601String(),
      if (location != null) 'location': location!.toJson(),
      // prefer 'readings' key per backend example; include raw_signal_data as alias
      'readings': rawSignalData.map((r) {
        if (r is SignalPoint) return r.toJson();
        if (r is Map<String, dynamic>) return r;
        return r;
      }).toList(),
      'raw_signal_data': rawSignalData.map((r) {
        if (r is SignalPoint) return r.toJson();
        if (r is Map<String, dynamic>) return r;
        return r;
      }).toList(),
      if (analysis != null) 'analysis': analysis!.toJson(),
    };
  }

  String toRawJson() => json.encode(toJson());
}
