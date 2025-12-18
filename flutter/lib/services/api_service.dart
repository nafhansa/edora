import 'dart:io';

import 'package:dio/dio.dart';
import '../models/reading_model.dart';

class SyncResponse {
  final bool success;
  final String? error;
  final int? statusCode;

  SyncResponse({required this.success, this.error, this.statusCode});
}

class ApiService {
  final Dio _dio;

  /// Default baseUrl points to Android emulator host loopback. Override in ctor.
  ApiService({String? baseUrl}) : _dio = Dio(BaseOptions(
    baseUrl: baseUrl ?? 'http://10.0.2.2:8080',
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 10),
    headers: {
      HttpHeaders.contentTypeHeader: 'application/json',
      HttpHeaders.acceptHeader: 'application/json',
    },
  ));

  /// Attempt to sync a reading to backend POST /api/v1/sync/reading
  /// Returns SyncResponse with detailed error info on failure.
  Future<SyncResponse> syncReading(Reading reading, {String? token}) async {
    try {
      final payload = reading.toJson();
      final headers = <String, dynamic>{};
      if (token != null && token.isNotEmpty) headers[HttpHeaders.authorizationHeader] = 'Bearer $token';

      final response = await _dio.post(
        '/api/v1/sync/reading',
        data: payload,
        options: Options(headers: headers),
      );

      if (response.statusCode != null && response.statusCode! >= 200 && response.statusCode! < 300) {
        return SyncResponse(success: true, statusCode: response.statusCode);
      }

      return SyncResponse(success: false, error: 'unexpected_status', statusCode: response.statusCode);
    } on DioException catch (e) {
      // Network, timeout, or server error - return descriptive error so caller can persist locally
      if (e.type == DioExceptionType.connectionTimeout || e.type == DioExceptionType.sendTimeout || e.type == DioExceptionType.receiveTimeout) {
        return SyncResponse(success: false, error: 'timeout', statusCode: null);
      }

      if (e.response != null) {
        final code = e.response?.statusCode;
        final msg = e.response?.data != null ? e.response?.data.toString() : e.message;
        return SyncResponse(success: false, error: msg, statusCode: code);
      }

      return SyncResponse(success: false, error: e.message, statusCode: null);
    } catch (e) {
      return SyncResponse(success: false, error: e.toString(), statusCode: null);
    }
  }
}
