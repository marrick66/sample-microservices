using System;

namespace postprocessing.Integration.Reactive
{
    /// <summary>
    /// Generic disposable that takes closures for managed and unmanaged 
    /// cleanup and runs them. Since Dispose() may be called concurrently,
    /// the closures should handle any synchronization needed.
    /// </summary>
    public class ActionDisposable : IDisposable
    {
        private Action _managed, _unmanaged;
        private bool _isDisposed = false;

        public ActionDisposable(Action Managed, Action Unmanaged = null)
        {
            _managed = Managed ??
                throw new ArgumentNullException(nameof(Managed));

            _unmanaged = Unmanaged;
        }

        /// <summary>
        /// Implement the standard dispose pattern for managed and
        /// unmanaged resources.
        /// </summary>
        public void Dispose()
        {
            Dispose(true);
        }

        protected void Dispose(bool IsDisposing)
        {
            if (_isDisposed)
                return;

            if (IsDisposing)
                _managed();

            _unmanaged?.Invoke();

            _isDisposed = true;

        }
    }
}
