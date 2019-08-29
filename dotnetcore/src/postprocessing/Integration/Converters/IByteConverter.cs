namespace postprocessing.Integration.Converters
{
    public interface IByteConverter
    {
        object FromBytes(byte[] Bytes);
        T FromBytes<T>(byte[] Bytes);

        byte[] ToBytes(object Item);
        byte[] ToBytes<T>(T Item);


    }
}
