﻿using Newtonsoft.Json;
using System.Text;

namespace postprocessing.Integration.Converters
{
    /// <summary>
    /// Basic converter class that gets the JSON string
    /// from the raw bytes, and deserializes it to an
    /// object instance.
    /// </summary>
    public class JsonByteConverter : IByteConverter
    {
        private Encoding _encoding;

        public JsonByteConverter(Encoding Encoder = null)
        {
            _encoding = Encoder ?? Encoding.UTF8;
        }

        public object FromBytes(byte[] Bytes)
        {
            var json = _encoding.GetString(Bytes);
            return JsonConvert.DeserializeObject(json);
        }

        public T FromBytes<T>(byte[] Bytes)
        {
            var json = _encoding.GetString(Bytes);
            return JsonConvert.DeserializeObject<T>(json);
        }

        public byte[] ToBytes(object Obj)
        {
            var json = JsonConvert.SerializeObject(Obj);
            return _encoding.GetBytes(json);
        }

        public byte[] ToBytes<T>(T Obj)
        {
            var json = JsonConvert.SerializeObject(Obj);
            return _encoding.GetBytes(json);
        }
    }
}
