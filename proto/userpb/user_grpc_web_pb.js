/**
 * @fileoverview gRPC-Web generated client stub for userpb
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.userpb = require('./user_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.userpb.UserServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.userpb.UserServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.userpb.User,
 *   !proto.userpb.UserResponse>}
 */
const methodDescriptor_UserService_UserTestCall = new grpc.web.MethodDescriptor(
  '/userpb.UserService/UserTestCall',
  grpc.web.MethodType.UNARY,
  proto.userpb.User,
  proto.userpb.UserResponse,
  /**
   * @param {!proto.userpb.User} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userpb.UserResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.userpb.User,
 *   !proto.userpb.UserResponse>}
 */
const methodInfo_UserService_UserTestCall = new grpc.web.AbstractClientBase.MethodInfo(
  proto.userpb.UserResponse,
  /**
   * @param {!proto.userpb.User} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.userpb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.userpb.User} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.userpb.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.userpb.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.userpb.UserServiceClient.prototype.userTestCall =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/userpb.UserService/UserTestCall',
      request,
      metadata || {},
      methodDescriptor_UserService_UserTestCall,
      callback);
};


/**
 * @param {!proto.userpb.User} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.userpb.UserResponse>}
 *     A native promise that resolves to the response
 */
proto.userpb.UserServicePromiseClient.prototype.userTestCall =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/userpb.UserService/UserTestCall',
      request,
      metadata || {},
      methodDescriptor_UserService_UserTestCall);
};


module.exports = proto.userpb;

