import asyncio

from grpclib.client import Channel

from ozonmp.lgc_location_api.v1.lgc_location_api_grpc import LgcLocationApiServiceStub
from ozonmp.lgc_location_api.v1.lgc_location_api_pb2 import DescribeLocationV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = LgcLocationApiServiceStub(channel)

        req = DescribeLocationV1Request(location_id=1)
        reply = await client.DescribeLocationV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())
