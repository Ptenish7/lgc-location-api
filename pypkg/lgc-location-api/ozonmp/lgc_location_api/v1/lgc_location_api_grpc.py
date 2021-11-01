# Generated by the Protocol Buffers compiler. DO NOT EDIT!
# source: ozonmp/lgc_location_api/v1/lgc_location_api.proto
# plugin: grpclib.plugin.main
import abc
import typing

import grpclib.const
import grpclib.client
if typing.TYPE_CHECKING:
    import grpclib.server

import validate.validate_pb2
import google.api.annotations_pb2
import google.protobuf.timestamp_pb2
import ozonmp.lgc_location_api.v1.lgc_location_api_pb2


class LgcLocationApiServiceBase(abc.ABC):

    @abc.abstractmethod
    async def CreateLocationV1(self, stream: 'grpclib.server.Stream[ozonmp.lgc_location_api.v1.lgc_location_api_pb2.CreateLocationV1Request, ozonmp.lgc_location_api.v1.lgc_location_api_pb2.CreateLocationV1Response]') -> None:
        pass

    @abc.abstractmethod
    async def DescribeLocationV1(self, stream: 'grpclib.server.Stream[ozonmp.lgc_location_api.v1.lgc_location_api_pb2.DescribeLocationV1Request, ozonmp.lgc_location_api.v1.lgc_location_api_pb2.DescribeLocationV1Response]') -> None:
        pass

    @abc.abstractmethod
    async def ListLocationsV1(self, stream: 'grpclib.server.Stream[ozonmp.lgc_location_api.v1.lgc_location_api_pb2.ListLocationsV1Request, ozonmp.lgc_location_api.v1.lgc_location_api_pb2.ListLocationsV1Response]') -> None:
        pass

    @abc.abstractmethod
    async def RemoveLocationV1(self, stream: 'grpclib.server.Stream[ozonmp.lgc_location_api.v1.lgc_location_api_pb2.RemoveLocationV1Request, ozonmp.lgc_location_api.v1.lgc_location_api_pb2.RemoveLocationV1Response]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/CreateLocationV1': grpclib.const.Handler(
                self.CreateLocationV1,
                grpclib.const.Cardinality.UNARY_UNARY,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.CreateLocationV1Request,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.CreateLocationV1Response,
            ),
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/DescribeLocationV1': grpclib.const.Handler(
                self.DescribeLocationV1,
                grpclib.const.Cardinality.UNARY_UNARY,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.DescribeLocationV1Request,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.DescribeLocationV1Response,
            ),
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/ListLocationsV1': grpclib.const.Handler(
                self.ListLocationsV1,
                grpclib.const.Cardinality.UNARY_UNARY,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.ListLocationsV1Request,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.ListLocationsV1Response,
            ),
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/RemoveLocationV1': grpclib.const.Handler(
                self.RemoveLocationV1,
                grpclib.const.Cardinality.UNARY_UNARY,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.RemoveLocationV1Request,
                ozonmp.lgc_location_api.v1.lgc_location_api_pb2.RemoveLocationV1Response,
            ),
        }


class LgcLocationApiServiceStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.CreateLocationV1 = grpclib.client.UnaryUnaryMethod(
            channel,
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/CreateLocationV1',
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.CreateLocationV1Request,
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.CreateLocationV1Response,
        )
        self.DescribeLocationV1 = grpclib.client.UnaryUnaryMethod(
            channel,
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/DescribeLocationV1',
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.DescribeLocationV1Request,
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.DescribeLocationV1Response,
        )
        self.ListLocationsV1 = grpclib.client.UnaryUnaryMethod(
            channel,
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/ListLocationsV1',
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.ListLocationsV1Request,
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.ListLocationsV1Response,
        )
        self.RemoveLocationV1 = grpclib.client.UnaryUnaryMethod(
            channel,
            '/ozonmp.lgc_location_api.v1.LgcLocationApiService/RemoveLocationV1',
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.RemoveLocationV1Request,
            ozonmp.lgc_location_api.v1.lgc_location_api_pb2.RemoveLocationV1Response,
        )